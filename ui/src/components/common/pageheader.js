import React, { useState } from 'react';
import { Button, PageHeader } from 'antd';
import { SyncOutlined, SettingOutlined } from '@ant-design/icons';
import { useDispatch, useSelector } from 'react-redux';
import { sinkToHost, sinkToME } from '../store/room/room.actions';
import SettingsModal from './SettingsModal';

function Pageheader() {
  const { name, isHost, controls } = useSelector(state => state.room)
  const [isModalVisible, setIsModalVisible] = useState(false);
  const showModal = () => { setIsModalVisible(true); };
  const handleOk = () => { setIsModalVisible(false); };
  const handleCancel = () => { setIsModalVisible(false); };
  const dispatch = useDispatch()

  const getActionButtons = () => {
    const buttons = [];
    if (!isHost) {
      buttons.push(
        <Button type="primary" icon={<SyncOutlined />} key="3" onClick={() => dispatch(sinkToHost())}>Sync to host</Button>,
      );
    }
    if (controls || isHost) {
      buttons.push(
        <Button type="primary" icon={<SyncOutlined />} key="2" onClick={() => dispatch(sinkToME())}>Sync everyone to me</Button>,
      );
    }
    if (isHost) {
      buttons.push(
        <Button type="primary" onClick={() => setIsModalVisible(true)} icon={<SettingOutlined />} key="1" />,
      );
    }

    return buttons;
  };

  return (
    <PageHeader
      ghost={false}
      onBack={() => { props.leave(); }}
      title={`Room: ${name}`}
            // subTitle={currentlyPlaying()}
      extra={getActionButtons()}
    >
      <SettingsModal isModalVisible={isModalVisible} showModal={showModal} handleOk={handleOk} handleCancel={handleCancel} />
    </PageHeader>
  );
}
export default (Pageheader);
