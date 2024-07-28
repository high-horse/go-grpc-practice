const grpc = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');
const SERVER = 'localhost:50051';
// Load the protobuf files
const protoFiles = ["../proto/news_models.proto", "../proto/news_service.proto"];
const packageDefinition = protoLoader.loadSync(protoFiles, {
    keepCase: true,
    longs: String,
    enums: String,
    defaults: true,
    oneofs: true
});
const protoDescriptor = grpc.loadPackageDefinition(packageDefinition);
const newsProto = protoDescriptor.news; // Assuming the package name in your .proto files is 'news'
const client = new newsProto.Newservice(SERVER, grpc.credentials.createInsecure());
function getNewsBulk() {
    const request = {}; // Empty request for NewsRequest
    client.GetNewsBulk(request, (error, response) => {
        if (error) {
            console.error("Error:", error);
            return;
        }
        console.log("Response from the server:");
        response.news.forEach((news) => {
            const source = news.source;
            console.log(`Source ID: ${source.id}, Source Name: ${source.name}`);
            console.log(`Author: ${news.author}, Title: ${news.title}, Description: ${news.description}, Published At: ${news.publishedAt}`);
        });
    });
}
// Start the client
getNewsBulk();
//# sourceMappingURL=index.js.map