<template>
  <div class="register-container">
    <h1 class="title">Register</h1>
    <form @submit.prevent="registerUser" class="register-form">
      <input
        v-model="username"
        type="text"
        placeholder="Username"
        required
        class="input"
      />
      <input
        v-model="email"
        type="email"
        placeholder="Email"
        required
        class="input"
      />
      <input
        v-model="password"
        type="password"
        placeholder="Password"
        required
        class="input"
      />
      <button type="submit" class="button">Register</button>
    </form>
    <p v-if="message" class="message">{{ message }}</p>
  </div>
</template>

<script>
export default {
  data() {
    return {
      username: "",
      email: "",
      password: "",
      message: "",
    };
  },
  methods: {
    async registerUser() {
      try {
        const response = await fetch("http://localhost:8080/register", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            username: this.username,
            email: this.email,
            password: this.password,
          }),
        });

        if (!response.ok) {
          throw new Error("Failed to register");
        }

        const data = await response.json();
        this.message = data.message;
        navigateTo("/login")
      } catch (error) {
        console.error(error);
        this.message = "Registration failed!";
      }
    },
  },
};
</script>

<style scoped>
.register-container {
  max-width: 400px;
  margin: 50px auto;
  padding: 20px;
  background-color: #f9f9f9;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  border-radius: 10px;
  text-align: center;
  font-family: Arial, sans-serif;
}

.title {
  font-size: 1.8rem;
  color: #333;
  margin-bottom: 20px;
}

.register-form {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.input {
  padding: 10px;
  font-size: 1rem;
  border: 1px solid #ccc;
  border-radius: 5px;
  transition: border-color 0.3s ease;
}

.input:focus {
  outline: none;
  border-color: #007bff;
}

.button {
  padding: 10px;
  font-size: 1rem;
  color: white;
  background-color: #007bff;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

.button:hover {
  background-color: #0056b3;
}

.message {
  margin-top: 15px;
  font-size: 0.9rem;
  color: green;
}
</style>
