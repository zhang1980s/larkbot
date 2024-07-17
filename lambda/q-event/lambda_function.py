import requests
import os
import json
import boto3
from datetime import datetime
from requests_auth_aws_sigv4 import AWSSigV4


aws_auth = AWSSigV4('q',region='us-east-1')
# translate client
trans_client = boto3.client('translate', region_name='us-east-1')
s3_client = boto3.client('s3', region_name='ap-northeast-1')
sm_client = boto3.client('secretsmanager', region_name='ap-northeast-1')
ddb_client = boto3.client('dynamodb', region_name='ap-northeast-1')
bedrock_runtime = boto3.client(service_name='bedrock-runtime', region_name='us-east-1')

def bedrock_translate(content):
    system_prompt = 'You are a highly skilled translator with expertise in many languages. Your task is to identify the language of the text user provides and directly translate it into English while preserving the meaning, tone, and nuance of the original text. Please maintain proper grammar, spelling, and punctuation in English. Do not try to understand the content, just put the result in <res></res>. Never talk to user starting with "I apologize", just give the translated text.'

    # The payload to be provided to Bedrock 
    native_request_payload = {
        "anthropic_version": "bedrock-2023-05-31",
        "max_tokens": 2000,
        "system": system_prompt,
        "messages": [
            {
                "role": "user",
                "content": [{"type": "text", "text": content}],
            }
        ],
        "temperature": 0
    }

    # The actual call to retrieve an answer from the model
    response = bedrock_runtime.invoke_model(
        body=json.dumps(native_request_payload),
        modelId='anthropic.claude-3-haiku-20240307-v1:0'
    )
    response_body = json.loads(response.get('body').read())
    return response_body


def get_secret(secret):
    get_secret_value_response = sm_client.get_secret_value(SecretId=secret)
    secret = get_secret_value_response['SecretString']
    return secret

ddb_respone = ddb_client.get_item(TableName=os.environ['CFG_TABLE'], Key={'key': {'S': os.environ['CFG_KEY']}})

def isEnglish(s):
    try:
        s.encode(encoding='utf-8').decode('ascii')
    except UnicodeDecodeError:
        return False
    else:
        return True


def translate(input_text,language):
    
    trans_response = trans_client.translate_text(
        Text=input_text,
        SourceLanguageCode='auto',
        TargetLanguageCode=language,
    )
    print(trans_response.get('TranslatedText'))
    return trans_response.get('TranslatedText')



def sendBackMessage(messageId,res):
    token_url = "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal"
    payload = {
        "app_id": get_secret(ddb_respone['Item']['app_id_arn']['S']),
        "app_secret": get_secret(ddb_respone['Item']['app_secret_arn']['S'])
    }
    headers = {'Content-type': 'application/json; charset=utf-8'}
    # Send the POST request
    response = requests.request(
        'POST',
        token_url,
        json=payload,
        headers=headers)
    print("token respose---", response)

    js = json.loads(response.text)
    tenant_access_token = js.get('tenant_access_token')
    print("tenant_access_token---", 'Bearer '+tenant_access_token)
    reply_url = "https://open.feishu.cn/open-apis/im/v1/messages/"+messageId+"/reply"

    headers = {
        'Content-type': 'application/json; charset=utf-8',
        'Authorization': 'Bearer '+tenant_access_token
    }
    
    content = json.dumps(res)
    content = content.replace('\n', '\\n')
    
    payload = {
        "content": content,
        "msg_type": "text",
        "reply_in_thread": True
    }
    print("rick "+json.dumps(payload))
    # Send the POST request
    response = requests.request(
        'POST',
        reply_url,
        json=payload,
        headers=headers)
    print(response.text)

    return response.text


def lambda_handler(event, context):
    print("incoming: " + json.dumps(event))
    res = ""

    text = json.loads(event['Records'][0]['body'])['content']
    print("text---", text)
    # use aws translate
    # if not isEnglish(text):
    #     text = translate(text,'en')

    text = bedrock_translate('"'+ text + '"' + ' -> EN').get('content')[0].get('text')
    if not text.find('<res>'):
        text = text.replace('<res>', '').replace('</res>', '')
    else:
        text = translate(text,'en')
    print("translated_text---", text)   
    res = res + requestAmazonQ(text)
    print("res---", res)

    # form json object
    content = {}
    content["text"] = res
    #send back to feishu
    #read from DDB
    sendBackMessage(event['Records'][0]['attributes']['MessageGroupId'],content)
    file_name = 'Q-'+ str(datetime.today().strftime('%Y-%m-%d-%H-%M-%S'))+'.txt'
    filepath = '/tmp/' + file_name
    
    with open(filepath, "xt") as f:
        f.write("Question: " + text + "\n"+ "Answer: " + res + "\n")
        f.close()
    
    
    s3_client.upload_file(filepath, os.environ['LOG_BUCKET'], datetime.today().strftime('%Y-%m-%d')+ "/"+file_name)
    os.remove(filepath)
    
    
    return res


# get response from amazon q
def requestAmazonQ(text):
    # get token
    url = "https://q.us-east-1.amazonaws.com/StartConversation"
    payload = {"source": "CONSOLE"}

    # Send the POST request
    response = requests.request('POST',
                                url,
                                json=payload,
                                auth=aws_auth
                                )
    print("token respose---", response)

    js = json.loads(response.text)
    conver_id = js.get('conversationId')
    conver_token = js.get('conversationToken')

    # get result
    url = "https://q.us-east-1.amazonaws.com/SendMessage"

    payload = {"source": "CONSOLE",
               "conversationId": conver_id,
               "utterance": text,
               "conversationToken": conver_token}

    # Send the POST request
    response = requests.request('POST',
                                url,
                                json=payload,
                                auth=aws_auth
                                )
    print(response.text)

    js = json.loads(response.text)
    ref = js['result']['content']['text']['references']
    ref_text = ''
    cnt = 1
    for i in ref:
        title = i['title']
        url = i['url']
        ref_text += str(cnt) + ". [{}]: {}".format(title, url) + '\n'
        cnt += 1
    res = json.loads(response.text)['result']['content']['text']['body']
    # use aws translate
    # translated_res = translate(res,'zh')
    translated_res = bedrock_translate('"'+ res + '"' + ' -> CN').get('content')[0].get('text')
    if not translated_res.find('<res>'):
        translated_res = translated_res.replace('<res>', '').replace('</res>', '')
    else:
        translated_res = translate(res,'zh')
    if ref_text != '':
        res = translated_res + '\n\n' + res + '\n\n' +  ref_text
    return res

