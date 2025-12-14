'use client';

const AUTH_TOKEN_COOKIE = 'auth_token';

/**
 * Set auth token in both localStorage and cookie (client-side)
 */
export function setAuthTokenCookie(token: string) {
  // Set in localStorage for client-side access
  if (typeof window !== 'undefined') {
    localStorage.setItem('auth_token', token);
    
    // Set in cookie for server-side access
    document.cookie = `${AUTH_TOKEN_COOKIE}=${token}; path=/; max-age=31536000; SameSite=Lax`;
  }
}

/**
 * Remove auth token from both localStorage and cookie (client-side)
 */
export function removeAuthTokenCookie() {
  if (typeof window !== 'undefined') {
    localStorage.removeItem('auth_token');
    localStorage.removeItem('auth_username');
    localStorage.removeItem('auth_user_id');
    
    // Remove cookie
    document.cookie = `${AUTH_TOKEN_COOKIE}=; path=/; max-age=0`;
  }
}
