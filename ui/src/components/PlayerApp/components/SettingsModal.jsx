import React, { useState, useEffect, useContext } from 'react';
import {
  Modal, Form, Button, Switch,
} from 'antd';

import {
  SettingOutlined,
} from '@ant-design/icons';

import { openNotificationWithIconKey } from '../../common/notification';
import { RoomContext } from '../../../context/RoomContext';

const layout = {
  labelCol: { span: 12 },
  wrapperCol: { span: 12 },
};

const SettingsModal = ({ isModalVisible, handleOk, handleCancel, queue} ) => {
  
  const [room, { updateQueue, updateSettings }] = useContext(RoomContext);
  const [controls, setControls] = useState(room.controls);
  const [auto_skip, setAutoSkip] = useState(room.auto_skip);

  useEffect(() => {
    setControls(room.controls);
    setAutoSkip(room.auto_skip);
  }, [room.controls, room.auto_skip]);

  const submitForm = () => {
    updateSettings(controls, auto_skip);
    handleOk();
  };
  const cancelForm = () => {
    setControls(room.controls);
    handleCancel();
  };

  const handleDarkMode = () => {
    openNotificationWithIconKey('warning', 'Only Dark Mode for you!', 'darkmode');
    const videoList = [...room.queue];
    videoList.unshift({ url: 'https://www.youtube.com/watch?v=dQw4w9WgXcQ', user: 'Watch2Gether' });
    updateQueue(videoList);
  };

  return (
    <Modal
      title="Room Settings"
      open={isModalVisible}
      onCancel={cancelForm}
      footer={[
        <Button key="back" onClick={cancelForm}>
          Cancel
        </Button>,
        <Button key="submit" type="primary" onClick={submitForm}>
          Submit
        </Button>,
      ]}
    >
      <Form {...layout} name="basic">
        <Form.Item
          label="Enable Player Sink To Me Button"
          name="controls"
        >
          <Switch checked={controls} onChange={() => (setControls(!controls))} />
        </Form.Item>
        <Form.Item
          label="Auto skip when everyone finishes"
          name="auto skip"
        >
          <Switch checked={auto_skip} onChange={() => (setAutoSkip(!auto_skip))} />
        </Form.Item>
        <Form.Item
          label="Dark Mode"
          name="controls"
        >
          <Switch defaultChecked checked onClick={handleDarkMode} />
        </Form.Item>
      </Form>
    </Modal>
  );
};

export const Settings = () => {
  const [isSettingsModalVisible, setIsSettingModalVisible] = useState(false);
  const showSettingsModal = () => { setIsSettingModalVisible(true); };
  const handleSettingsOk = () => { setIsSettingModalVisible(false); };
  const handleSettingsCancel = () => { setIsSettingModalVisible(false); };

  return (
    <div style={{ width: '100%' }}>
      <Button style={{ width: '100%', padding: '0px 8px', height: '33px' }} type="primary" onClick={() => setIsSettingModalVisible(true)} icon={<SettingOutlined />} key="1" />
      <SettingsModal isModalVisible={isSettingsModalVisible} showModal={showSettingsModal} handleOk={handleSettingsOk} handleCancel={handleSettingsCancel} />
    </div>
  );
};

