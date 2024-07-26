<template>
    <div class="flex justify-between items-center">
      <b>Categories</b>
      <UButton
        icon="i-heroicons-pencil-square"
        size="sm"
        color="primary"
        variant="solid"
        :trailing="false"
        @click="$emit('new-content')"
      />
    </div>
    <hr />
  </template>