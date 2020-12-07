import { notification } from 'antd';

export const openNotificationWithIcon = (type, message) => {
    notification[type]({
      key: "stopspam",
      message: message,
      position: "bottomRight"
    });
  };
