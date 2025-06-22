<script setup lang="ts">
import { RouterView, RouterLink } from "vue-router";
</script>

<template>
  <div>
    <!-- Simple navigation bar -->
    <nav class="nav-bar">
      <h1>Juuten Game</h1>
      <div class="nav-links">
        <RouterLink to="/">Home</RouterLink>
        <RouterLink to="/signup">Sign Up</RouterLink>
        <RouterLink to="/rooms">Rooms</RouterLink>
        <RouterLink to="/play">Play</RouterLink>
        <RouterLink to="/api-test" class="api-test-link">🛠️ API Test</RouterLink>
      </div>
    </nav>
    
    <main>
      <router-view />
    </main>
  </div>
</template>

<style>
.nav-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 2rem;
  background: #f8f9fa;
  border-bottom: 1px solid #e9ecef;
  margin-bottom: 1rem;
}

.nav-bar h1 {
  margin: 0;
  color: #2c3e50;
  font-size: 1.5rem;
}

.nav-links {
  display: flex;
  gap: 1rem;
}

.nav-links a {
  text-decoration: none;
  color: #007bff;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.nav-links a:hover {
  background-color: #e7f1ff;
}

.nav-links a.router-link-active {
  background-color: #007bff;
  color: white;
}

.api-test-link {
  background-color: #28a745 !important;
  color: white !important;
  font-weight: bold;
}

.api-test-link:hover {
  background-color: #218838 !important;
}

main {
  min-height: calc(100vh - 100px);
}

@media (max-width: 768px) {
  .nav-bar {
    flex-direction: column;
    gap: 1rem;
  }
  
  .nav-links {
    flex-wrap: wrap;
    justify-content: center;
  }
}
</style>
