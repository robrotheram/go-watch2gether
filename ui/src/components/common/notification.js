import { notification } from 'antd';

export const openNotificationWithIcon = (type, message) => {
  openNotificationWithIconKey(type, message, "stopSPAM");
  };

export const openNotificationWithIconKey = (type, message, key) => {
  notification[type]({
    key: key,
    message: message,
    placement: "topRight"
  });
};
