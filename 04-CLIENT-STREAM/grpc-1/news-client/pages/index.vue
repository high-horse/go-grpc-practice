<template>
    <div class="flex flex-wrap">
        <div class="w-full md:w-2/12">
            <MenuCategoryComponent 
            @category-clicked="clickedCategory"
            @new-content="generateNewContent"
            />
        </div>
        <div class="w-full md:w-2/12">
            <MenuListComponent />
        </div>
        <div class="w-full md:w-8/12">
            <MenuDetailComponent />    
        </div>
    </div>
</template>



<script setup lang="ts">
const { getNewsBulk } = useNewsService()

const clickedCategory = async (key: string) => {
    try {
        await getNewsBulk ()
    } catch (e) {
        console.error('Failed to fetch news:', e)
    }
}

const generateNewContent = () => {
    // todo generate new content
    console.log("generate new content");
}


onMounted(() => {
    console.log("mounted to the index page");
    
})

</script>