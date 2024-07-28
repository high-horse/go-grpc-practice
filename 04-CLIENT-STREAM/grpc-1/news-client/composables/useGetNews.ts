import type { UseFetchOptions } from "#app";

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

export const useGetNews = (
  endpoint: string,
  method: string = 'GET',
  payload: any = null
) => {
  const BASE_URL = "http://localhost:8000/";
  const url = BASE_URL + endpoint;

  const news: Ref<NewsItem[] | null> = ref(null);
  const error: Ref<Error | null> = ref(null);
  const isLoading = ref(false);

  const fetchNews = async () => {
    isLoading.value = true;
    error.value = null;

    try {
      const response = await fetch(url, {
        method,
        headers: {
          'Content-Type': 'application/json',
        },
        body: payload ? JSON.stringify(payload) : null,
      });

      if (!response.ok) {
        throw new Error('Network response was not ok');
      }

      const data: NewsResponse = await response.json();
      news.value = data.news;
    } catch (e) {
      error.value = e as Error;
    } finally {
      isLoading.value = false;
    }
  };

  return {
    news,
    error,
    isLoading,
    fetchNews
  }
};
