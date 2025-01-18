document.addEventListener('DOMContentLoaded', () => {
    const registerForm = document.getElementById('registerForm');
    const alertBox = document.getElementById('alert');
  
    function showAlert(message, isSuccess = false) {
      alertBox.textContent = message;
      alertBox.style.display = 'block';
      alertBox.className = 'alert ' + (isSuccess ? 'alert-success' : 'alert-error');
    }
  
    if (registerForm) {
      registerForm.addEventListener('submit', async (e) => {
        e.preventDefault();
  
        const username = document.getElementById('username').value;
        const email = document.getElementById('email').value;
        const password = document.getElementById('password').value;
        const fullName = document.getElementById('fullName').value;
        const mobileNumber = document.getElementById('mobileNumber').value;
        const birthDate = document.getElementById('birthDate').value;
        const gender = document.getElementById('gender').value;
  
        try {
          const resp = await fetch('/register', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
              username,
              email,
              password,
              full_name: fullName,
              mobile_number: mobileNumber,
              birth_date: birthDate,
              gender
            }),
          });
          const data = await resp.json();
  
          if (resp.ok) {
            showAlert('Registration success: ' + data.message, true);
          } else {
            showAlert('Registration error: ' + (data.error || 'Unknown error'), false);
          }
        } catch (error) {
          showAlert('Error: ' + error.message, false);
        }
      });
    }
  });
  