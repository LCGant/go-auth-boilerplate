import win32com.client as win32
import pythoncom
from flask import Flask, request, jsonify # type: ignore

app = Flask(__name__)

def send_email(destination_email, subject, body):
    """Função para enviar e-mail com Outlook"""
    try:
        pythoncom.CoInitialize()
        
        outlook = win32.Dispatch('outlook.application')
        mail = outlook.CreateItem(0)

        mail.To = destination_email
        mail.Subject = subject
        mail.Body = body

        mail.Send()
        return True
    except Exception as e:
        print(f"Erro ao enviar e-mail: {e}")
        return False
    finally:
        pythoncom.CoUninitialize()

@app.route('/send-email', methods=['POST'])
def handle_email():
    """Rota para receber JSON e enviar o e-mail"""
    data = request.get_json()
    if not data or 'email' not in data or 'subject' not in data or 'message' not in data:
        return jsonify({"success": False, "error": "Invalid JSON format"}), 400

    destination_email = data['email']
    subject = data['subject']
    message = data['message']

    if send_email(destination_email, subject, message):
        return jsonify({"success": True, "message": "Email sent successfully"})
    else:
        return jsonify({"success": False, "message": "Failed to send email"}), 500

if __name__ == '__main__':
    app.run(port=5000)
