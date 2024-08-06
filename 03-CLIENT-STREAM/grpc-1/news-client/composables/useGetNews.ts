import type { UseFetchOptions } from "#app";
import { log } from "@grpc/grpc-js/build/src/logging";
import type { NewsItem, NewsResponse } from "../types";

export const useGetNews = (
  endpoint: Ref<string>,
  method: string = "GET",
  payload: any = null,
) => {
  const BASE_URL = "http://localhost:8000/";

  var url = computed(() => BASE_URL + endpoint.value);

  const news: Ref<NewsItem[] | null> = ref(null);
  const error: Ref<Error | null> = ref(null);
  const isLoading = ref(false);

  const fetchNews = async () => {
    isLoading.value = true;
    error.value = null;

    try {
      const response = await fetch(url.value, {
        method,
        headers: {
          "Content-Type": "application/json",
        },
        body: payload ? JSON.stringify(payload) : null,
      });

      if (!response.ok) {
        throw new Error("Network response was not ok");
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
    fetchNews,
  };
};
