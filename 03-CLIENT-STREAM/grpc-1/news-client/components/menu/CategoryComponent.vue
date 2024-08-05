<template>
    <div class="border-div">
        <MenuCategoryHeadingComponent @new-content="emitNewContent()" />
       <div>
            <ul>

                <li v-for="(value, key) in options"
                :key="key"
                :class="{active: key == activeOption}"
                @click="setActiveOption(key)"
                class="cursor-pointer"
                >
                    {{ value }}
                </li>
            </ul>
       </div>
    </div>
    
</template>

<script setup lang="ts">

const emit = defineEmits<{
    (event: 'category-clicked', key: string): void;
    (event: 'new-content'): void;
}>();

const options = {
    "fresh-news": "Recent News",
    "db-news" : "Fetch DB News",
    "get_created_news" : "Get Created News",
    "create_news":"Create News",
};

const activeOption = ref<string | null>(null);


const setActiveOption = (key: string) => {
    activeOption.value = key;
    if (key) {
        emit('category-clicked', key);
    }
};

const emitNewContent = () => {
    setActiveOption("")
    emit('new-content')
}

</script>

<style scoped>
.active {
  background-color: #1376f8;
  font-weight: bold;
}
.cursor-pointer {
  cursor: pointer;
}
.border-div {
    border: solid wheat;
    border-radius: 10px;
    height: 100vh;
    padding: 10px;
}
</style>