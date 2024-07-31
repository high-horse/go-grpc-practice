<template>
    <div>
        <b> Heading</b>
        <UTextarea
            color="primary"
            variant="outline"
            placeholder="Heading..."
            model-value=""
        />
        <b>Author</b>
        <UInput color="primary" variant="outline" placeholder="Author..." />
        <b>Source</b>
        <UInput color="primary" variant="outline" placeholder="Source..." />
        <b>Description:</b>
        <UTextarea 
            ref="textareaRef" 
            color="primary"
            autoresize
            placeholder="News Description..." 
            :rows="5"
            :maxrows="20"
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
</template>

<script setup lang="ts">

</script>

