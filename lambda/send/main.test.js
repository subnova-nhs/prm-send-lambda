const ehrExtract = require("./main");
const given = require('./given')
const AWS = require("aws-sdk-mock");

describe("When passing in invalid payloads", () => {

  test("a bad request response is returned, for an empty string", async () => {
    let event = {"body": "" };
    let result = await ehrExtract.handler(event);
    expect(result.statusCode).toBe(400);
  });

  test("a bad request response is returned, for an empty event", async () => {
    let event = {};
    let result = await ehrExtract.handler(event);
    expect(result.statusCode).toBe(400);
  });

  test("a bad request response is returned, for invalid xml", async () => {
    let event = {"body": "blah" };
    let result = await ehrExtract.handler(event);
    expect(result.statusCode).toBe(400);
  });
});

describe("ERROR responses", () => {
  let result;

  beforeAll(async () => {
    AWS.mock('DynamoDB.DocumentClient', 'put', function (params, callback){
      callback(null, Promise.reject('Oops!'));
    });
    let event = {"body": given.tpp_sample_encodedXml };
    result = await ehrExtract.handler(event);
  });

  test("That if there is an error when saving the data, it generates a ERROR response", async () => {
    expect(result.body).toBe("{\"status\":\"ERROR\"}");
  });

  afterAll(() => {
    AWS.restore('DynamoDB.DocumentClient');
  });
});

describe("ACCEPTED responses", () => {
  let result;
  let dynamoDbPutCallCount = 0;

  beforeAll(async () => {
    AWS.mock('DynamoDB.DocumentClient', 'put', function(params, callback) {
      dynamoDbPutCallCount++;
      if (params.Item.PROCESS_PAYLOAD !== given.tpp_sample_json) {
        throw "Payload does not match expected";
      } 
      callback(null, {});
    });
    let event = {"body": given.tpp_sample_encodedXml};
    result = await ehrExtract.handler(event);
  });

  test("returns a successful response", async () => {
    expect(result.statusCode).toBe(200);
  });

  test("it should store data", async () => {
    expect(dynamoDbPutCallCount).toBe(1);
  });

  test("That if there is no error, it generates an ACCEPTED response", async () => {
    var jsonBody = JSON.parse(result.body);
    expect(jsonBody.status).toBe("ACCEPTED");
  });

  afterAll(() => {
    AWS.restore('DynamoDB.DocumentClient');
  });
});
