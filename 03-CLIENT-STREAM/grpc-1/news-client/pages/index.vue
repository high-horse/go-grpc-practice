<template>
    <div class="flex flex-wrap">
      <div class="w-full md:w-2/12">
        <MenuCategoryComponent 
        @new-content="generateNewContent"
          @category-clicked="clickedCategory"
        />
      </div>
      <div class="w-full md:w-2/12">
        <MenuListComponent 
          :news="news"
          :error="error"
          :isLoading="isLoading"
          :fetchNews="fetchNews"
          @news-clicked="newsSelected"
        />
      </div>
      <div class="w-full md:w-8/12">
        <MenuDetailComponent :selectedNews="selectedNews"/>    
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import { ref, watchEffect, onMounted } from 'vue';
  import type { NewsItem } from '~/types';
  
  const { news, error, isLoading, fetchNews } = useGetNews('news', 'GET', null);
  const selectedNews = ref<NewsItem | null>(null);
  
  watchEffect(() => {
    if (news.value) {
      console.log('News updated:', news.value);
    }
    if (error.value) {
      console.error('Error occurred:', error.value);
    }
  });
  
  const clickedCategory = async (key: string) => {
    try {
      console.log(`Category clicked: ${key}`);
      selectedNews.value = null;
      await fetchNews(); // Trigger fetching news only when a category is clicked
    } catch (e) {
      console.error('Failed to fetch news:', e);
    }
  };
  
  const newsSelected = (title: string) => {
    console.log("Selected news: ", title);
    const newsItem = news.value?.find(item => item.title === title) || null;
    selectedNews.value = newsItem;
  };
  
  const generateNewContent = () => {
    console.log("Generating new content");
    
    // Create a new "UNTITLED" news item
    const untitledNews: NewsItem = {
      source: { id: "", name: "" },
      author: "",
      title: "UNTITLED",
      description: "",
      publishedAt: "",
    };
    news.value = [];
    // Add the new "UNTITLED" item to the news list
    if (news.value) {
      news.value = [untitledNews];
    } else {
      news.value = [untitledNews];
    }
  
    // Set the selected news to the new "UNTITLED" item
    selectedNews.value = untitledNews;
  };
  
  onMounted(() => {
    console.log("Mounted to the index page");
  });
</script>
  