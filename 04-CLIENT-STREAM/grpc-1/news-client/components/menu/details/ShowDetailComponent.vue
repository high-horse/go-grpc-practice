<!-- MenuDetailsShowDetailComponent.vue -->
<template>
    <div>
      <b>
        <h1>{{ selectedNews.title }}</h1>
      </b>
      <p>
        <i>Author: {{ selectedNews.author || "unknown" }},&nbsp;&nbsp;&nbsp;&nbsp;</i>
        <i>source: {{ selectedNews.source.name || "unknown" }}</i>
      </p>
      
      <p>Description:</p>
      <UTextarea 
        ref="textareaRef" 
        autoresize
        :disabled="isDisabled"
        placeholder="News Description..." 
        v-model="selectedNews.description" 
        :rows="20"
        :maxrows="20"
      />
          
      <p class="flex justify-end"><i>{{ selectedNews.publishedAt }}</i></p>
      <div>
        
        <div class="flex justify-end" style="">
          <UButton
            icon="i-heroicons-pencil-square"
            size="sm"
            color="primary"
            variant="solid"
            label="Edit"
            :trailing="false"
            @click="enableEdit"
            class="mr-2"
          />
          <UButton
            icon="i-heroicons-pencil-square"
            size="sm"
            color="primary"
            variant="solid"
            label="Save"
            :trailing="true"
          />
        </div>
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
    // import { ref, watch } from 'vue';
    import type { NewsItem } from '~/types';
    
    const props = defineProps<{
        selectedNews: NewsItem;
    }>();
    
    const isDisabled = ref(true);
    
    const enableEdit = () => {
        isDisabled.value = false;
    };
    
    watch(
        () => props.selectedNews,
        () => {
        isDisabled.value = true;
        }
    );
  </script>
  