package email

const AppointmentConfirmationTemplate = `<!DOCTYPE html>
<html>
<head><meta charset="utf-8"></head>
<body style="font-family:Arial,sans-serif;max-width:600px;margin:0 auto;padding:20px;color:#333">
  <div style="background:#0e7490;padding:20px;border-radius:8px 8px 0 0;text-align:center">
    <h1 style="color:#fff;margin:0">Arnela Gabinete</h1>
  </div>
  <div style="border:1px solid #e5e7eb;border-top:none;padding:30px;border-radius:0 0 8px 8px">
    <h2 style="color:#0e7490">Tu cita ha sido confirmada</h2>
    <p>Hola <strong>{{.ClientName}}</strong>,</p>
    <p>Tu cita ha sido confirmada con los siguientes detalles:</p>
    <table style="width:100%;border-collapse:collapse;margin:20px 0">
      <tr><td style="padding:8px;border-bottom:1px solid #e5e7eb;font-weight:bold">Servicio</td><td style="padding:8px;border-bottom:1px solid #e5e7eb">{{.Title}}</td></tr>
      <tr><td style="padding:8px;border-bottom:1px solid #e5e7eb;font-weight:bold">Fecha</td><td style="padding:8px;border-bottom:1px solid #e5e7eb">{{.Date}}</td></tr>
      <tr><td style="padding:8px;border-bottom:1px solid #e5e7eb;font-weight:bold">Hora</td><td style="padding:8px;border-bottom:1px solid #e5e7eb">{{.Time}}</td></tr>
      <tr><td style="padding:8px;border-bottom:1px solid #e5e7eb;font-weight:bold">Profesional</td><td style="padding:8px;border-bottom:1px solid #e5e7eb">{{.EmployeeName}}</td></tr>
      {{if .Room}}<tr><td style="padding:8px;border-bottom:1px solid #e5e7eb;font-weight:bold">Sala</td><td style="padding:8px;border-bottom:1px solid #e5e7eb">{{.Room}}</td></tr>{{end}}
    </table>
    <p style="background:#f0fdfa;padding:12px;border-radius:6px;border-left:4px solid #0e7490">
      Si necesitas cancelar o modificar tu cita, contacta con nosotros al <strong>611 749 043</strong> o responde a este correo.
    </p>
    <hr style="border:none;border-top:1px solid #e5e7eb;margin:20px 0">
    <p style="font-size:12px;color:#9ca3af">Arnela Gabinete · C/ García Barbón 30, Vigo · 611 749 043</p>
  </div>
</body>
</html>`

const AppointmentCancellationTemplate = `<!DOCTYPE html>
<html>
<head><meta charset="utf-8"></head>
<body style="font-family:Arial,sans-serif;max-width:600px;margin:0 auto;padding:20px;color:#333">
  <div style="background:#0e7490;padding:20px;border-radius:8px 8px 0 0;text-align:center">
    <h1 style="color:#fff;margin:0">Arnela Gabinete</h1>
  </div>
  <div style="border:1px solid #e5e7eb;border-top:none;padding:30px;border-radius:0 0 8px 8px">
    <h2 style="color:#dc2626">Cita cancelada</h2>
    <p>Hola <strong>{{.ClientName}}</strong>,</p>
    <p>Tu cita del <strong>{{.Date}}</strong> a las <strong>{{.Time}}</strong> ({{.Title}}) ha sido cancelada.</p>
    {{if .Reason}}<p><strong>Motivo:</strong> {{.Reason}}</p>{{end}}
    <p>Si deseas reprogramar, puedes hacerlo desde nuestra web o llamando al <strong>611 749 043</strong>.</p>
    <hr style="border:none;border-top:1px solid #e5e7eb;margin:20px 0">
    <p style="font-size:12px;color:#9ca3af">Arnela Gabinete · C/ García Barbón 30, Vigo · 611 749 043</p>
  </div>
</body>
</html>`

type AppointmentEmailData struct {
	ClientName   string
	Title        string
	Date         string
	Time         string
	EmployeeName string
	Room         string
	Reason       string
}
