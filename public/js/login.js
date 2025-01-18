document.addEventListener('DOMContentLoaded', () => {
  const alertBox = document.getElementById('alert');
  const loginForm = document.getElementById('loginForm');

  function showAlert(message, isSuccess = false) {
    alertBox.textContent = message;
    alertBox.style.display = 'block';
    alertBox.className = 'alert ' + (isSuccess ? 'alert-success' : 'alert-error');
  }

  if (loginForm) {
    loginForm.addEventListener('submit', async (e) => {
      e.preventDefault();
      const username = document.getElementById('username').value;
      const password = document.getElementById('password').value;

      try {
        const resp = await fetch('/login', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ username, password }),
        });
        const data = await resp.json();

        if (resp.ok) {
          showAlert('Login success: ' + data.message, true);
        } else {
          showAlert('Login error: ' + (data.error || 'Unknown error'), false);
        }
      } catch (error) {
        showAlert('Error: ' + error.message, false);
      }
    });
  }

  const forgotLink = document.getElementById('forgotLink');
  const forgotContainer = document.getElementById('forgotContainer');
  const forgotEmail = document.getElementById('forgotEmail');
  const forgotBtn = document.getElementById('forgotBtn');

  if (forgotLink) {
    forgotLink.addEventListener('click', (e) => {
      e.preventDefault();
      forgotContainer.style.display = (forgotContainer.style.display === 'none') ? 'block' : 'none';
    });
  }

  if (forgotBtn) {
    forgotBtn.addEventListener('click', async () => {
      const email = forgotEmail.value;
      if (!email) {
        showAlert('Por favor, insira o e-mail.', false);
        return;
      }
      try {
        const resp = await fetch('/forgot-password', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ email }),
        });
        const data = await resp.json();

        if (resp.ok) {
          showAlert('Reset: ' + data.message, true);
        } else {
          showAlert('Erro ao solicitar reset: ' + (data.error || 'Unknown error'), false);
        }
      } catch (err) {
        showAlert('Erro: ' + err.message, false);
      }
    });
  }
});
