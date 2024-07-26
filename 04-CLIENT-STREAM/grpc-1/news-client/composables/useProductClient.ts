import grpc from '@grpc/grpc-js';
import protoLoader from '@grpc/proto-loader';

const SERVER = 'localhost:50051';
// current path: E:\go-raty\grpc\grpc-1\04-CLIENT-STREAM\grpc-1\news-client\composables\grpcClient.ts
// path to proto files:
// E:\go-raty\grpc\grpc-1\04-CLIENT-STREAM\grpc-1\proto\news_models.proto
// E:\go-raty\grpc\grpc-1\04-CLIENT-STREAM\grpc-1\proto\news_service.proto
const protoFiles = ["../proto/news_models.proto", "../proto/news_service.proto"];

const packageDefinition = protoLoader.loadSync(protoFiles, {
    keepCase: true,
    longs: String,
    enums: String,
    defaults: true,
    oneofs: true
});

const protoDescriptor = grpc.loadPackageDefinition(packageDefinition);
const newsProto = protoDescriptor.news as any; // Adjust as needed

const client = new newsProto.Newservice(SERVER, grpc.credentials.createInsecure());

async function useGetNewsBulk(): Promise<any> {
    return new Promise((resolve, reject) => {
        const req = {};

        client.GetNewsBulk(req, (error: grpc.ServiceError | null, response: any) => {
            if (error) {
                reject(error);
                return;
            }
            resolve(response);
        });
    });
}

export { useGetNewsBulk };
