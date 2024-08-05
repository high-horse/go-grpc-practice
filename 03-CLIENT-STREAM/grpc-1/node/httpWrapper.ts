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
  url: string;
    publishedAt: string;
}

interface NewsResponse {
    news: NewsItem[];
}

const STATUS_CODES = {
    OK: 200,
    NO_CONTENT: 204,
    NOT_FOUND: 404,
    INTERNAL_SERVER_ERROR: 500,
};

/* GRPC SERVER CALLS START */
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

function GetFreshNews() :Promise<NewsResponse> {
  return new  Promise((resolve, reject) => {
    const request = {};
    client.GetFreshNews(request, (error: grpc.ServiceError | null, response: NewsResponse) => {
      if (error ) {
        reject(error);
        return;
      }
      resolve(response);
    })
  })
}

function getDBNews() :Promise<NewsResponse> {
  return new  Promise((resolve, reject) => {
    const request = {};
    client.GetDBNews(request, (error: grpc.ServiceError | null, response: NewsResponse) => {
      if (error ) {
        reject(error);
        return;
      }
      resolve(response);
    })
  })
}
/* GRPC SERVER CALLS ENDS */

/*HTTP MIDDLEWARE STARTS */
const logHttp = (req: http.IncomingMessage, res: http.ServerResponse, newt :() => void) => {
  const { method, url } = req;
  console.log(`[${new Date().toISOString()}] ${method} ${url}`);
  newt();
}
/*HTTP MIDDLEWARE ENDS */

/* HTTP HANDLERS START */
const requestHandler = async (req: http.IncomingMessage, res: http.ServerResponse) => {
  res.setHeader('Access-Control-Allow-Origin', '*');
  res.setHeader('Access-Control-Allow-Methods', 'GET, POST, OPTIONS');
  res.setHeader('Access-Control-Allow-Headers', 'Content-Type');

  if (req.method === 'OPTIONS') {
    res.writeHead(STATUS_CODES.NO_CONTENT);
    res.end();
    return;
  }

  try { 
    let newsResponse: NewsResponse | null = null;
    if (req.url === "/fresh-news" && req.method === "GET") {
      newsResponse = await GetFreshNews();
    } else if (req.url === "/db-news" && req.method === "GET") {
      newsResponse = await getDBNews();
    } else {
      res.writeHead(STATUS_CODES.NOT_FOUND, { 'Content-type': 'application/json' });
      res.end(JSON.stringify({ error: "Not Found" }));
      return;
    }
    res.writeHead(STATUS_CODES.OK, { 'Content-type': 'application/json' });
    res.end(JSON.stringify(newsResponse))
  } catch(error) {
    console.log("error:", error);
    res.writeHead(STATUS_CODES.INTERNAL_SERVER_ERROR, { 'Content-type': 'application/json' });
    res.end(JSON.stringify({ error: "Internal Server Error" }));
  }
};

const server = http.createServer((req, res) => {
  logHttp(req, res, () => {
    requestHandler(req, res);
  })
})

server.listen(HTTP_PORT, () => {
  console.log(`HTTP server running on http://localhost:${HTTP_PORT}`);
  console.log("/fresh-news      -> Fetch Fresh News");
  console.log("/db-news         -> Fetch News from DB");
  
});
