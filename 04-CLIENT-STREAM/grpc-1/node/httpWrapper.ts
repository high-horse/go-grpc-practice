import * as grpc from '@grpc/grpc-js';
import * as protoLoader from '@grpc/proto-loader';
import * as http from 'http';

const SERVER = 'localhost:50051';
const HTTP_PORT = 8000;

const protoFiles = ["../proto/news_models.proto", "../proto/news_service.proto"];
const packageDefinition = protoLoader.loadSync(protoFiles, {
    keepCase: true,
    longs: String,
    enums: String,
    defaults: true,
    oneofs: true
});

const protoDescriptor = grpc.loadPackageDefinition(packageDefinition) as any;
const newsProto = protoDescriptor.news as any;
const client = new newsProto.Newservice(SERVER, grpc.credentials.createInsecure());

interface NewsSource {
    id: string;
    name: string;
}

interface NewsItem {
    source: NewsSource;
    author: string;
    title: string;
    description: string;
    publishedAt: string;
}

interface NewsResponse {
    news: NewsItem[];
}

function getNewsBulk(): Promise<NewsResponse> {
    return new Promise((resolve, reject) => {
        const request = {};
        client.GetNewsBulk(request, (error: grpc.ServiceError | null, response: NewsResponse) => {
            if (error) {
                reject(error);
                return;
            }
            resolve(response);
        });
    });
}

const server = http.createServer(async (req, res) => {
    // res.setHeader('Access-Control-Allow-Origin', '*');
    // res.setHeader('Access-Control-Allow-Methods', 'GET, POST, OPTIONS');
    // res.setHeader('Access-Control-Allow-Headers', 'Content-Type');
    

    const allowedOrigins = ['http://localhost:3000'];
    const origin = req.headers.origin as string;

    if (allowedOrigins.includes(origin)) {
        res.setHeader('Access-Control-Allow-Origin', origin);
        res.setHeader('Access-Control-Allow-Methods', 'GET, POST, OPTIONS');
        res.setHeader('Access-Control-Allow-Headers', 'Content-Type');
    }

    if (req.method === 'OPTIONS') {
        res.writeHead(204);
        res.end();
        return;
    }

    if (req.url === '/news' && req.method === 'GET') {
        try {
            const newsResponse = await getNewsBulk();
            res.writeHead(200, { 'Content-Type': 'application/json' });
            res.end(JSON.stringify(newsResponse));
        } catch (error) {
            console.error("Error:", error);
            res.writeHead(500, { 'Content-Type': 'application/json' });
            res.end(JSON.stringify({ error: 'Internal Server Error' }));
        }
    } else {
        res.writeHead(404, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: 'Not Found' }));
    }
});

server.listen(HTTP_PORT, () => {
    console.log(`HTTP server running on http://localhost:${HTTP_PORT}`);
});
