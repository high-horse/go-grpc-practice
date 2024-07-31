<template>
  <form @submit.prevent="handleSubmit">
    <div>
        
        <b> Heading</b>
        <UTextarea
            color="primary"
            variant="outline"
            placeholder="Heading..."
            v-model="formData.title"
        />
        <b>Author</b>
        <UInput 
          color="primary" 
          variant="outline" 
          placeholder="Author..." 
          v-model="formData.author"
        />
        <b>Source</b>
        <UInput 
          color="primary" 
          variant="outline" 
          placeholder="Source..." 
          v-model="formData.source.name"
        />
        <b>Description:</b>
        <UTextarea 
            ref="textareaRef" 
            color="primary"
            autoresize
            placeholder="News Description..." 
            :rows="5"
            :maxrows="20"
            v-model="formData.description"
        />
        

    </div>
    <div>
        <div class="flex justify-end" style="padding: 10px;">
          <UButton
            icon="i-heroicons-pencil-square"
            size="sm"
            color="primary"
            variant="solid"
            label="Discard"
            :trailing="false"
            class="mr-2"
            @click="discardNews"
          />
          <UButton
            icon="i-heroicons-pencil-square"
            size="sm"
            color="primary"
            variant="solid"
            label="Save"
            type="submit"
            :trailing="true"

          />
        </div>
    </div>

  </form>
</template>

<script setup lang="ts">
import type { NewsItem } from '~/types';


const formData = ref<NewsItem >({
  source: { id: '', name: '' },
  title: '',
  author: '',
  description: '',
  publishedAt: '',
});

// const emits = defineEmits<{
//   (event: 'submit-newnews', submitData: NewsItem): void;
// }>();

const discardNews = () => {
  formData.value = {
    source: {id:'', name: ''},
    title: "",
    author: "",
    description: '',
    publishedAt: '',
  }
}
const handleSubmit = () => {
  console.log("FORMDATA", formData.value);
}
</script>

