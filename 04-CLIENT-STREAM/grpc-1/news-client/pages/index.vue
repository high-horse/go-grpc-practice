<template>
    <div class="flex flex-wrap">
        <div class="w-full md:w-2/12">
            <MenuCategoryComponent @item-clicked="clickedCategory"/>
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

const clickedCategory = (key: string) => {
    console.log("item clicked: ", key)
}
</script>