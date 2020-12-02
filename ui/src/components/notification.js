import { notification } from 'antd';

export const openNotificationWithIcon = (type, message) => {
    notification[type]({
      message: message,
      position: "bottomRight"
    });
  };
