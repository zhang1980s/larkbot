import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb'
import * as apigateway from 'aws-cdk-lib/aws-apigateway'
import * as iam from 'aws-cdk-lib/aws-iam'
import * as secretsmanager from 'aws-cdk-lib/aws-secretsmanager'
import * as events from 'aws-cdk-lib/aws-events'
import * as targets from 'aws-cdk-lib/aws-events-targets'


export class LarkbotAppStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);


    ///////////////////////////////////////////////////////////////////////
    // Define AppID and AppSecret as cfn input parameters
    ///////////////////////////////////////////////////////////////////////

    const appID = new cdk.CfnParameter(this, 'AppID', {
      type: 'String',
      description: 'The AppID of larkbot app',
      noEcho: true,
      default: 'cli_xxxxxxxxxxxxxxxx',
    })

    const appSecret = new cdk.CfnParameter(this, 'AppSecret',{
      type: 'String',
      description: 'The Secret ID of larkbot app',
      noEcho: true,
      default: 'XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX',
    })

    const caseLanguage = new cdk.CfnParameter(this, 'CaseLanguage',{
      type: 'String',
      description: 'Case Language queue. Should be in "zh", "ja", "ko", "en" ',
      noEcho: false,
      allowedValues: ["zh","ja","ko","en"],
      default: 'zh'
    })

    const configKey = new cdk.CfnParameter(this, 'ConfigKey', {
      type: 'String',
      description: 'The default config profile',
      noEcho: false,
      default: 'LarkBotProfile-0'
    })

    const userWhitelist = new cdk.CfnParameter(this, 'UserWhitelist',{
      type: 'String',
      description: 'Enable user white list function',
      noEcho: false,
      allowedValues: ["true","false"],
      default: 'false'
    })

    const supportRegion = new cdk.CfnParameter(this, 'SupportRegion', {
      type: 'String',
      description: 'The default support region',
      noEcho: false,
      allowedValues: ['en','cn'],
      default: 'en'
    })
    // const enableRefresh = new cdk.CfnParameter(this, 'EnableRefresh', {
    //   type:'String',
    //   description: 'Enable Refresh rule, disable by default',
    //   noEcho: false,
    //   allowedValues: ['true', 'false'],
    //   default:'false'
    // })

    const refreshInterval = new cdk.CfnParameter(this, 'RefreshInterval', {
      type: 'Number',
      description: 'Case refresh interval (in munute)',
      noEcho: false,
      default: 10
    })

    

    ///////////////////////////////////////////////////////////////////////
    // Define Secrets for AppID and AppSecret
    ///////////////////////////////////////////////////////////////////////


    const AppIDSecret = new secretsmanager.Secret(this, 'AppIDSecret', {
      description: 'The Secret to store the value of App ID',
      secretStringValue: cdk.SecretValue.cfnParameter(appID),
    
    })

    const AppSecretSecret = new secretsmanager.Secret(this, 'AppSecretSecret', {
      description: 'The Secret to store the value of App Secret',
      secretStringValue: cdk.SecretValue.cfnParameter(appSecret),
    })
    

    ///////////////////////////////////////////////////////////////////////
    // Define DDB tables 
    ///////////////////////////////////////////////////////////////////////

    const auditTable = new dynamodb.Table(this, 'audit', {
      partitionKey: {name: 'key', type: dynamodb.AttributeType.STRING },
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST
    })

    const botCasesTable = new dynamodb.Table(this, 'bot_cases', {
      partitionKey: {name: 'pk', type: dynamodb.AttributeType.STRING },
      sortKey: {name: 'sk', type: dynamodb.AttributeType.STRING},
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
    })


    botCasesTable.addGlobalSecondaryIndex(
      {
      indexName: 'card_msg_id-index',
      partitionKey: {
        name: 'card_msg_id',
        type: dynamodb.AttributeType.STRING,
      },
      projectionType: dynamodb.ProjectionType.ALL,
      }
    );

    botCasesTable.addGlobalSecondaryIndex(
      {
        indexName: 'status-type-index',
        partitionKey: {
          name: 'status',
          type: dynamodb.AttributeType.STRING,
        },
        sortKey: {
          name: 'type',
          type: dynamodb.AttributeType.STRING,
        },
        projectionType: dynamodb.ProjectionType.ALL,
      }
    );

    const botConfigTable = new dynamodb.Table(this, 'bot_config', {
      partitionKey: {name: 'key', type: dynamodb.AttributeType.STRING },
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST
    })




    ///////////////////////////////////////////////////////////////////////
    // Define lambda functions with alias and version
    ///////////////////////////////////////////////////////////////////////

    // Define msgEvent handler
    const msgEventFunction = new lambda.Function(this,'larkbot-msg-event', {
      runtime: lambda.Runtime.PROVIDED_AL2,
      architecture: lambda.Architecture.ARM_64,
      handler: 'bootstrap',
      code: lambda.Code.fromAsset('lambda/msg-event'),
      // tracing: lambda.Tracing.ACTIVE,
      timeout: cdk.Duration.minutes(1),
      environment: {
        AUDIT_TABLE: auditTable.tableName,
        CASES_TABLE: botCasesTable.tableName,
        CFG_TABLE: botConfigTable.tableName,
        CFG_KEY: configKey.valueAsString,
        CASE_LANGUAGE: caseLanguage.valueAsString,
        ENABLE_USER_WHITELIST: userWhitelist.valueAsString,
        SUPPORT_REGION: supportRegion.valueAsString,
       }
    } );


    const msgEventVersion = msgEventFunction.currentVersion;

    const msgEventAlias = new lambda.Alias(this, 'msg-event-prod', {
      aliasName: 'Prod',
      version: msgEventVersion,
    });

    // Attch the policy document that allow to access Secret ARN of the AppID and AppSecret

    AppIDSecret.grantRead(msgEventAlias)
    AppSecretSecret.grantRead(msgEventAlias)


    // Attach the policy document that allow to assume the support role in others accounts to the lambda function's role
        msgEventAlias.addToRolePolicy(new iam.PolicyStatement(
          {
            sid: 'AllowToAssumeToRoleWithSupportAPIAccess',
            effect: iam.Effect.ALLOW,
            actions: ['sts:AssumeRole'],
            resources: ['arn:aws:iam::*:role/FeishuSupportCaseApiAll*']
          }
        ))

    // Grant RW access of audit table to larkbot function 

    auditTable.grantReadWriteData(msgEventAlias)
    botCasesTable.grantReadWriteData(msgEventAlias)
    botConfigTable.grantReadWriteData(msgEventAlias)


    ///////////////////////////////////////////////////////////////////////
    // Define the Rest APIs for message and content card 
    ///////////////////////////////////////////////////////////////////////

    const msgEventApi = new apigateway.LambdaRestApi(this, 'msgEventapi', {
      handler: msgEventAlias,
      proxy: false,
    })

    const eventMessages = msgEventApi.root.addResource('messages');

    eventMessages.addMethod(
      'POST', 
      new apigateway.LambdaIntegration(msgEventAlias, {
      proxy: false,
      integrationResponses: [
        {
          statusCode: '200',
          responseTemplates: {
            'application/json': '',
          }
        },
      ],
    }),
    {
      methodResponses: [
        {
          statusCode: "200",
          responseModels: {
            "application/json": apigateway.Model.EMPTY_MODEL
          }
        },
        {
          statusCode: "400",
          responseModels: {
            "application/json": apigateway.Model.ERROR_MODEL
          }
        },
        {
          statusCode: "500",
          responseModels: {
            "application/json": apigateway.Model.ERROR_MODEL
          }
        }
      ]
    })

    ///////////////////////////////////////////////////////////////////////
    // Define Eventbridge rule
    ///////////////////////////////////////////////////////////////////////


    const refreshEventRule = new events.Rule(this, 'refreshCaseRule', {
      schedule: events.Schedule.rate(cdk.Duration.minutes(refreshInterval.valueAsNumber)),
      description: `Refresh case update every ${refreshInterval.valueAsString} minutes`,
      // enabled: Boolean(enableRefresh.valueAsString === 'true')
      enabled: false
    })

    refreshEventRule.addTarget(new targets.LambdaFunction(msgEventAlias, {
      event: events.RuleTargetInput.fromObject(
        {
          schema: "2.0",
          event: {
            message: {
              message_type: "fresh_comment"
            }
          }
        }
      )
    }))


    ///////////////////////////////////////////////////////////////////////
    // Output the roleArn of msgEvent
    ///////////////////////////////////////////////////////////////////////

    const msgEventAliasRole = msgEventAlias.role

    new cdk.CfnOutput(this, 'msgEventRoleArn', {
      value: msgEventAliasRole!.roleArn ,
      description: 'The arn of msgEventfunction',
    });

    }

  }
