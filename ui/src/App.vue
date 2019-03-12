<template>
  <section class="section">
    <div id="app" class="container">
      <Landing />
    </div>
  </section>
</template>

<script>
import Landing from './components/Landing.vue'

export default {
  name: 'app',
  components: {
    Landing
  }
}
</script>