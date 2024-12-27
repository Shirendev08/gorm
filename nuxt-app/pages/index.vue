<template>
  <div>
    <h1>Welcome to the Home Page</h1>
    <p v-if="username">Hello, {{ username }}!</p>
    <button @click="logout">Logout</button>
  </div>
</template>

<script>
export default {
  data() {
    return {
      username: '', // This will store the username
    };
  },
  mounted() {
    // Get the username from cookies
    this.username = this.$cookies.get('username');
  },
  methods: {
    logout() {
      // Clear the cookies
      this.$cookies.remove('authToken');
      this.$cookies.remove('username');

      // Redirect to the login page
      this.$router.push('/login');
    },
  },
};
</script>
