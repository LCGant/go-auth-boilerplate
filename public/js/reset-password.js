document.addEventListener('DOMContentLoaded', () => {
    const urlParams = new URLSearchParams(window.location.search);
    const token = urlParams.get('token');
    const tokenInput = document.getElementById('token');
    tokenInput.value = token;

    const resetForm = document.getElementById('resetForm');
    resetForm.addEventListener('submit', async (e) => {
      e.preventDefault();
      const newPassword = document.getElementById('newPassword').value;

      const resp = await fetch('/reset-password', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ token, newPassword })
      });
      const data = await resp.json();

      if (resp.ok) {
        alert('Senha redefinida com sucesso!');
      } else {
        alert('Erro: ' + data.error);
      }
    });
  });