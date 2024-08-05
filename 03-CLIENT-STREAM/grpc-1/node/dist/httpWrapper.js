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
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
Object.defineProperty(exports, "__esModule", { value: true });
const grpc = __importStar(require("@grpc/grpc-js"));
const protoLoader = __importStar(require("@grpc/proto-loader"));
const http = __importStar(require("http"));
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
const protoDescriptor = grpc.loadPackageDefinition(packageDefinition);
const newsProto = protoDescriptor.news;
const client = new newsProto.Newservice(SERVER, grpc.credentials.createInsecure());
const STATUS_CODES = {
    OK: 200,
    NO_CONTENT: 204,
    NOT_FOUND: 404,
    INTERNAL_SERVER_ERROR: 500,
};
/* GRPC SERVER CALLS START */
function getNewsBulk() {
    return new Promise((resolve, reject) => {
        const request = {};
        client.GetNewsBulk(request, (error, response) => {
            if (error) {
                reject(error);
                return;
            }
            resolve(response);
        });
    });
}
function GetFreshNews() {
    return new Promise((resolve, reject) => {
        const request = {};
        client.GetFreshNews(request, (error, response) => {
            if (error) {
                reject(error);
                return;
            }
            resolve(response);
        });
    });
}
function getDBNews() {
    return new Promise((resolve, reject) => {
        const request = {};
        client.GetDBNews(request, (error, response) => {
            if (error) {
                reject(error);
                return;
            }
            resolve(response);
        });
    });
}
/* GRPC SERVER CALLS ENDS */
/*HTTP MIDDLEWARE STARTS */
const logHttp = (req, res, newt) => {
    const { method, url } = req;
    console.log(`[${new Date().toISOString()}] ${method} ${url}`);
    newt();
};
/*HTTP MIDDLEWARE ENDS */
const requestHandler = (req, res) => __awaiter(void 0, void 0, void 0, function* () {
    res.setHeader('Access-Control-Allow-Origin', '*');
    res.setHeader('Access-Control-Allow-Methods', 'GET, POST, OPTIONS');
    res.setHeader('Access-Control-Allow-Headers', 'Content-Type');
    if (req.method === 'OPTIONS') {
        res.writeHead(STATUS_CODES.NO_CONTENT);
        res.end();
        return;
    }
    try {
        let newsResponse = null;
        if (req.url === "/fresh-news" && req.method === "GET") {
            newsResponse = yield GetFreshNews();
        }
        else if (req.url === "/db-news" && req.method === "GET") {
            newsResponse = yield getDBNews();
        }
        else {
            res.writeHead(STATUS_CODES.NOT_FOUND, { 'Content-type': 'application/json' });
            res.end(JSON.stringify({ error: "Not Found" }));
            return;
        }
        res.writeHead(STATUS_CODES.OK, { 'Content-type': 'application/json' });
        res.end(JSON.stringify(newsResponse));
    }
    catch (error) {
        console.log("error:", error);
        res.writeHead(STATUS_CODES.INTERNAL_SERVER_ERROR, { 'Content-type': 'application/json' });
        res.end(JSON.stringify({ error: "Internal Server Error" }));
    }
});
const server = http.createServer((req, res) => {
    logHttp(req, res, () => {
        requestHandler(req, res);
    });
});
server.listen(HTTP_PORT, () => {
    console.log(`HTTP server running on http://localhost:${HTTP_PORT}`);
    console.log("/fresh-news      -> Fetch Fresh News");
    console.log("/db-news         -> Fetch News from DB");
});
//# sourceMappingURL=httpWrapper.js.map