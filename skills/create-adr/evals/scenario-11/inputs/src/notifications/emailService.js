const nodemailer = require('nodemailer');
const logger = require('../logger');

const transporter = nodemailer.createTransport({
  host: process.env.SMTP_HOST,
  port: parseInt(process.env.SMTP_PORT, 10),
  secure: true,
  auth: {
    user: process.env.SMTP_USER,
    pass: process.env.SMTP_PASS,
  },
});

// Sends order confirmation email synchronously — called directly from the
// POST /orders request handler. Slow SMTP responses delay the HTTP response.
async function sendOrderConfirmation(order) {
  const info = await transporter.sendMail({
    from: '"NotifyCore Billing" <billing@notifycore.io>',
    to: order.customerEmail,
    subject: `Order Confirmed: ${order.id}`,
    html: buildOrderConfirmationHtml(order),
  });
  logger.info('Order confirmation sent', { messageId: info.messageId, orderId: order.id });
  return info;
}

// Sends payment receipt synchronously — called from the payment webhook handler.
async function sendPaymentReceipt(payment) {
  const info = await transporter.sendMail({
    from: '"NotifyCore Billing" <billing@notifycore.io>',
    to: payment.customerEmail,
    subject: `Payment Receipt: ${payment.id}`,
    html: buildPaymentReceiptHtml(payment),
  });
  logger.info('Payment receipt sent', { messageId: info.messageId, paymentId: payment.id });
  return info;
}

// Sends password reset email synchronously — called from the auth route.
async function sendPasswordReset(user, resetToken) {
  const info = await transporter.sendMail({
    from: '"NotifyCore" <no-reply@notifycore.io>',
    to: user.email,
    subject: 'Password Reset Request',
    html: buildPasswordResetHtml(user, resetToken),
  });
  logger.info('Password reset sent', { messageId: info.messageId, userId: user.id });
  return info;
}

function buildOrderConfirmationHtml(order) {
  return `<p>Your order <strong>${order.id}</strong> has been confirmed.</p>`;
}

function buildPaymentReceiptHtml(payment) {
  return `<p>Payment <strong>${payment.id}</strong> received. Amount: $${payment.amount}.</p>`;
}

function buildPasswordResetHtml(user, token) {
  return `<p>Click <a href="${process.env.APP_URL}/reset?token=${token}">here</a> to reset your password.</p>`;
}

module.exports = { sendOrderConfirmation, sendPaymentReceipt, sendPasswordReset };
