export default defineNuxtRouteMiddleware((to, from) => {
    const authToken = useCookie('authToken');
  
    // If not authenticated and trying to access any page other than login, redirect to /login
    if (!authToken) {
      return navigateTo('/login');
    }
  
  
  
    // Allow navigation for all other cases
  });
  