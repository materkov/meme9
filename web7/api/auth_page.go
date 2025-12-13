package api

import (
	"fmt"
	"html"
	"net/http"
)

func (a *API) authPageHandler(w http.ResponseWriter, r *http.Request) {
	// Get query parameter for tab (login or register)
	tab := r.URL.Query().Get("tab")
	if tab != "register" {
		tab = "login"
	}

	// Check for error messages from query params (for redirects after failed auth)
	errorMsg := r.URL.Query().Get("error")
	usernameError := r.URL.Query().Get("usernameError")
	credentialsError := r.URL.Query().Get("credentialsError")

	// Build error HTML (escape to prevent XSS)
	usernameErrorHTML := ""
	if usernameError != "" {
		usernameErrorHTML = fmt.Sprintf(`<div class="fieldError">%s</div>`, html.EscapeString(usernameError))
	}

	credentialsErrorHTML := ""
	if credentialsError != "" {
		credentialsErrorHTML = fmt.Sprintf(`<div class="fieldError">%s</div>`, html.EscapeString(credentialsError))
	}

	errorHTML := ""
	if errorMsg != "" {
		errorHTML = fmt.Sprintf(`<div class="error">%s</div>`, html.EscapeString(errorMsg))
	}

	// Determine active tab classes
	loginTabClass := ""
	registerTabClass := ""
	if tab == "login" {
		loginTabClass = "active"
	} else {
		registerTabClass = "active"
	}

	// Submit button text
	submitText := "Login"
	if tab == "register" {
		submitText = "Register"
	}

	// Password autocomplete
	passwordAutocomplete := "current-password"
	if tab == "register" {
		passwordAutocomplete = "new-password"
	}

	htmlContent := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Login - meme9</title>
  <style>
    body {
      min-height: 100vh;
      display: flex;
      align-items: center;
      justify-content: center;
      background: #f5f5f5;
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
      margin: 0;
    }
    .container {
      background: white;
      border-radius: 8px;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
      width: 100%s;
      max-width: 400px;
      padding: 2rem;
    }
    .tabs {
      display: flex;
      border-bottom: 1px solid #e0e0e0;
      margin-bottom: 1.5rem;
    }
    .tab {
      flex: 1;
      padding: 0.75rem;
      background: none;
      border: none;
      cursor: pointer;
      font-size: 1rem;
      color: #666;
      border-bottom: 2px solid transparent;
      transition: all 0.2s;
      text-decoration: none;
      display: block;
      text-align: center;
    }
    .tab:hover {
      color: #333;
    }
    .tab.active {
      color: #333;
      border-bottom-color: #333;
      font-weight: 500;
    }
    .form {
      display: flex;
      flex-direction: column;
      gap: 1.5rem;
    }
    .field {
      display: flex;
      flex-direction: column;
      gap: 0.5rem;
    }
    .field label {
      font-size: 0.9rem;
      font-weight: 500;
      color: #333;
    }
    .field input {
      padding: 0.75rem;
      border: 1px solid #ddd;
      border-radius: 4px;
      font-size: 1rem;
      transition: border-color 0.2s;
    }
    .field input:focus {
      outline: none;
      border-color: #333;
    }
    .inputError {
      border-color: #dc3545 !important;
    }
    .fieldError {
      color: #dc3545;
      font-size: 0.875rem;
      margin-top: 0.25rem;
    }
    .error {
      padding: 0.75rem;
      background: #fee;
      border: 1px solid #fcc;
      border-radius: 4px;
      color: #c33;
      font-size: 0.9rem;
    }
    .submit {
      padding: 0.75rem;
      background: #333;
      color: white;
      border: none;
      border-radius: 4px;
      font-size: 1rem;
      font-weight: 500;
      cursor: pointer;
      transition: background 0.2s;
    }
    .submit:hover {
      background: #555;
    }
  </style>
</head>
<body>
  <div class="container">
    <div class="tabs">
      <a href="/?tab=login" class="tab %s">Login</a>
      <a href="/?tab=register" class="tab %s">Register</a>
    </div>

    <form id="authForm" class="form">
      <div class="field">
        <label for="username">Username</label>
        <input
          id="username"
          name="username"
          type="text"
          required
          autocomplete="username"
          %s
        />
        %s
      </div>

      <div class="field">
        <label for="password">Password</label>
        <input
          id="password"
          name="password"
          type="password"
          required
          autocomplete="%s"
          %s
        />
        %s
      </div>

      %s

      <button type="submit" class="submit">%s</button>
    </form>
  </div>

  <script>
    // Handle form submission with fetch API
    document.getElementById('authForm').addEventListener('submit', async function(e) {
      e.preventDefault();
      
      const form = e.target;
      const username = document.getElementById('username').value;
      const password = document.getElementById('password').value;
      const submitButton = form.querySelector('.submit');
      const originalButtonText = submitButton.textContent;
      
      // Disable form during submission
      submitButton.disabled = true;
      submitButton.textContent = 'Loading...';
      
      // Clear previous errors
      document.querySelectorAll('.fieldError, .error').forEach(el => el.remove());
      document.querySelectorAll('.inputError').forEach(el => el.classList.remove('inputError'));
      
      const data = {
        username: username,
        password: password
      };

      // Determine endpoint based on current tab
      const currentTab = window.location.search.includes('tab=register') ? 'register' : 'login';
      const apiEndpoint = currentTab === 'login' ? '/api/login' : '/api/register';

      try {
        const response = await fetch(apiEndpoint, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(data)
        });

        if (response.ok) {
          const result = await response.json();
          // Store token and user info in localStorage
          localStorage.setItem('auth_token', result.token);
          localStorage.setItem('auth_user_id', result.user_id);
          localStorage.setItem('auth_username', result.username);
          // Set cookie for server-side auth checks
          document.cookie = 'auth_token=' + result.token + '; path=/; max-age=86400'; // 24 hours
          // Redirect to feed
          window.location.href = '/feed';
        } else {
          const error = await response.json();
          // Show error inline without page reload
          const usernameField = document.getElementById('username');
          const passwordField = document.getElementById('password');
          const usernameFieldDiv = usernameField.closest('.field');
          const passwordFieldDiv = passwordField.closest('.field');
          
          // Remove existing error messages
          usernameFieldDiv.querySelectorAll('.fieldError').forEach(el => el.remove());
          passwordFieldDiv.querySelectorAll('.fieldError').forEach(el => el.remove());
          document.querySelectorAll('.error').forEach(el => el.remove());
          
          // Remove error classes
          usernameField.classList.remove('inputError');
          passwordField.classList.remove('inputError');
          
          if (error.error_code === 'username_exists') {
            usernameField.classList.add('inputError');
            const errorDiv = document.createElement('div');
            errorDiv.className = 'fieldError';
            errorDiv.textContent = 'Username already exists';
            usernameFieldDiv.appendChild(errorDiv);
          } else if (error.error_code === 'invalid_credentials') {
            passwordField.classList.add('inputError');
            const errorDiv = document.createElement('div');
            errorDiv.className = 'fieldError';
            errorDiv.textContent = 'Invalid username or password';
            passwordFieldDiv.appendChild(errorDiv);
          } else {
            const errorDiv = document.createElement('div');
            errorDiv.className = 'error';
            errorDiv.textContent = error.error_details || error.error_code || 'An error occurred';
            form.insertBefore(errorDiv, form.querySelector('.submit'));
          }
          
          submitButton.disabled = false;
          submitButton.textContent = originalButtonText;
        }
      } catch (err) {
        // Show network error inline
        const form = document.getElementById('authForm');
        document.querySelectorAll('.error').forEach(el => el.remove());
        const errorDiv = document.createElement('div');
        errorDiv.className = 'error';
        errorDiv.textContent = 'Network error. Please try again.';
        form.insertBefore(errorDiv, form.querySelector('.submit'));
        
        submitButton.disabled = false;
        submitButton.textContent = originalButtonText;
      }
    });
  </script>
</body>
</html>`,
		"%",
		loginTabClass,
		registerTabClass,
		func() string {
			if usernameError != "" {
				return `class="inputError"`
			}
			return ""
		}(),
		usernameErrorHTML,
		passwordAutocomplete,
		func() string {
			if credentialsError != "" {
				return `class="inputError"`
			}
			return ""
		}(),
		credentialsErrorHTML,
		errorHTML,
		submitText,
	)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(htmlContent))
}
