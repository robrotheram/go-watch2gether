import React, { useState, useEffect } from 'react';
import { Modal, Input, Button } from 'antd';

import { connect } from 'react-redux';
import axios from 'axios';
import {
  CopyOutlined,
  ShareAltOutlined,
} from '@ant-design/icons';
import { openNotificationWithIconKey } from '../../common/notification';
import { BASE_URL } from '../../../context/config';


const ShareModal = ({isModalVisible, handleCancel}) => {

  const cancelForm = () => {
    handleCancel();
  };
  const [botid, setBot] = useState('');
  useEffect(() => {
    axios.get(`${BASE_URL}/config`).then((res) => {
      setBot(res.data.bot);
    });
  }, []);
  const inviteBotUrl = (bot) => `https://discord.com/oauth2/authorize?client_id=${bot}&scope=bot`;
  return (
    <Modal
      title="Share Room"
      open={isModalVisible}
      onCancel={cancelForm}
      footer={[
        <Button key="back" onClick={cancelForm}>
          Cancel
        </Button>,
      ]}
    >
      <Input.Group compact>
        <Input disabled style={{ width: '90%' }} defaultValue={window.location.href} />
        <Button
          style={{ width: '10%' }}
          type="primary"
          icon={<CopyOutlined />}
          onClick={() => {
            navigator.clipboard.writeText(window.location.href);
            openNotificationWithIconKey('success', 'Link coppied');
          }}
        />
      </Input.Group>
      <br />
      {botid !== '' ? (
        <Button
          target="_blank"
          href={inviteBotUrl(botid)}
          type="primary"
          style={{
            // padding: '0px 20px',
            backgroundColor: '#7289da',
            border: 'none',
            width: '100%',
          }}
        >
          Invite bot
        </Button>
      )
        : null}
    </Modal>
  );
};

const Share = () => {
  const [isShareModalVisible, setIsShareModalVisible] = useState(false);
  const showShareModal = () => { setIsShareModalVisible(true); };
  const handleShareOk = () => { setIsShareModalVisible(false); };
  const handleShareCancel = () => { setIsShareModalVisible(false); };

  return (
    <div style={{ width: '100%' }}>
      <Button style={{ width: '100%', padding: '0px 8px', height: '33px' }} type="primary" icon={<ShareAltOutlined />} onClick={() => setIsShareModalVisible(true)} />
      <ShareModal isModalVisible={isShareModalVisible} showModal={showShareModal} handleOk={handleShareOk} handleCancel={handleShareCancel} />
    </div>
  );
};
export default Share;
