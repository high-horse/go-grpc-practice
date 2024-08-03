// types.ts

export interface NewsSource {
    id: string;
    name: string;
  }
  
  export interface NewsItem {
    source: NewsSource;
    author: string;
    title: string;
    description: string;
    publishedAt: string;
  }
  
  export interface NewsResponse {
    news: NewsItem[];
  }
  