
client := aws.client('lambda')
client.invoke({FunctionName: 'example1'})['Payload'] | decode('base64') | json.unmarshal
