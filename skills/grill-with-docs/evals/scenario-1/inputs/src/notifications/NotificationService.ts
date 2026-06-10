import { sendEmail } from './adapters/email';
import { sendSMS } from './adapters/sms';
import { sendPush } from './adapters/push';

export class NotificationService {
  async notify(userId: string, type: 'email' | 'sms' | 'push', message: string) {
    if (type === 'email') await sendEmail(userId, message);
    if (type === 'sms') await sendSMS(userId, message);
    if (type === 'push') await sendPush(userId, message);
  }

  // Called synchronously from OrderService after order placement
  async notifyOrderPlaced(orderId: string, userId: string) {
    await this.notify(userId, 'email', `Order ${orderId} confirmed.`);
    await this.notify(userId, 'sms', `Order ${orderId} confirmed.`);
    await this.notify(userId, 'push', `Order ${orderId} confirmed.`);
  }
}
