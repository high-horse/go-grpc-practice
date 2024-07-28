"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || function (mod) {
    if (mod && mod.__esModule) return mod;
    var result = {};
    if (mod != null) for (var k in mod) if (k !== "default" && Object.prototype.hasOwnProperty.call(mod, k)) __createBinding(result, mod, k);
    __setModuleDefault(result, mod);
    return result;
};
Object.defineProperty(exports, "__esModule", { value: true });
const grpc = __importStar(require("@grpc/grpc-js"));
const protoLoader = __importStar(require("@grpc/proto-loader"));
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