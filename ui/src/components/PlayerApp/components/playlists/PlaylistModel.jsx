import React, { useState, useEffect } from 'react';
import {
  Modal, Button, Form, Input,
} from 'antd';
import {
  EditOutlined, PlusOutlined, DeleteOutlined,
} from '@ant-design/icons';
import { connect } from 'react-redux';
import { EditableTable } from './tables/editTable';

import { createPlaylists, updatePlaylists } from '../../../../store/playlists/playlists.actions';
import { createVideoItem, validURL } from '../../../../store/video';
import { updateQueue } from '../../../../store/room/room.actions';
import { openNotificationWithIcon } from '../../../common/notification';
import { ViewableTable } from './tables/viewTable';

const CREATE = 'c';
const UPDATE = 'u';

const PlaylistModel = ({
  visible, setVisible, data, title, room, queue, user, updateQueue, createPlaylists, updatePlaylists,
}) => {
  const [confirmLoading, setConfirmLoading] = useState(false);
  const [mode, setMode] = useState('VIEW');
  const [datastore, setDatastore] = useState([]);
  const [selected, setSelected] = useState([]);
  const [updateType, setType] = useState(CREATE);
  const [form] = Form.useForm();

  useEffect(() => {
    setMode('VIEW');
    if (title === 'Create new Playlist') {
      setMode('EDIT');
    }
  }, [visible, title]);

  useEffect(() => {
    if (data !== undefined) {
      form.setFieldsValue({ name: data.name });
      setType(UPDATE);
      setDatastore(data.videos.map((v) => {
        if (v.id === '') {
          v.id = ID();
          return v;
        }
        return v;
      }));
    } else {
      form.setFieldsValue({ name: '' });
      setType(CREATE);
      setDatastore([]);
    }
  }, [form, data]);

  const handleOk = () => {
    setConfirmLoading(true);
    savePlaylist();
  };

  const handleCancel = () => {
    console.log('Clicked cancel button');
    setVisible(false);
  };
  var ID = () => `_${Math.random().toString(36).substr(2, 9)}`;

  const addToPlaylist = () => {
    if (selected.length === 0) {
      setVisible(false);
      return;
    }
    const selectedVideos = datastore.filter((item) => selected.indexOf(item.id) !== -1);
    const newqueue = [...queue, ...selectedVideos];
    updateQueue(newqueue);
    setVisible(false);
  };

  const addrecord = () => {
    const data = {
      id: ID(),
      key: ID(),
      url: 'enter valid video url',
      order: datastore.length + 1,
    };
    setMode('EDIT');
    setDatastore((datastore) => [...datastore, data]);
  };

  const savePlaylist = async () => {
    const valid = datastore.every((item) => validURL(item.url));
    if (!valid) {
      setConfirmLoading(false);
      openNotificationWithIcon('error', 'Invalid URL');
      return;
    }
    if (form.getFieldsValue('name').name.length < 3) {
      setConfirmLoading(false);
      openNotificationWithIcon('error', 'Invalid: Playlist name needs to be greater then 3 characters');
      return;
    }

    const ds = await Promise.all(datastore.map(async (video) => { video = await createVideoItem(video.url, user); return video; }));

    if (updateType === CREATE) {
      const playlist = {
        name: form.getFieldsValue('name').name,
        username: '',
        videos: ds,
        room,
      };
      createPlaylists(room, playlist);
    } else {
      data.name = form.getFieldsValue('name').name;
      data.videos = ds;
      updatePlaylists(room, data);
    }

    setVisible(false);
    setConfirmLoading(false);
  };
  console.log('playlistModel', user, room);

  const modalTitle = (
    <div style={{ display: 'inline-flex', width: 'calc( 100% - 20px )' }}>
      <span style={{ marginRight: '10px' }}>
        {mode === 'EDIT' || mode === 'SORT'
          ? <Button disabled={datastore.length === 0} style={{ width: '90px' }} type="primary" icon={<PlusOutlined />} onClick={() => setMode('VIEW')}>View</Button>
          : <Button style={{ width: '90px' }} type="primary" icon={<EditOutlined />} onClick={() => setMode('EDIT')}>Edit</Button>}
      </span>
      <span style={{ width: '100%' }}>
        {mode === 'EDIT' || mode === 'SORT' ? (
          <Form form={form}>
            <Form.Item name="name" label="Edit name:" rules={[{ required: true }]} style={{ marginBottom: '0px' }}>
              <Input />
            </Form.Item>
          </Form>
        )
          : <p style={{ paddingTop: '5px', marginBottom: '0px' }}>{title}</p>}
      </span>
    </div>
  );

  const EditModeButton = () => {
    return (
      <Button icon={<DeleteOutlined />} key="edit" onClick={() => setMode('EDIT')} style={{ float: 'left' }}>
        Remove
      </Button>
    );
  };

  const modalFooter = () => {
    if (mode === 'EDIT' || mode === 'SORT') {
      return [
        EditModeButton(),
        <Button key="add" onClick={addrecord} style={{ float: 'left' }}>
          Add new Video
        </Button>,
        <Button key="back" onClick={handleCancel}>
          Cancel
        </Button>,
        <Button key="submit" type="primary" loading={confirmLoading} onClick={handleOk}>
          Submit
        </Button>,
      ];
    }
    return [
      <Button key="back" onClick={handleCancel}>
        Cancel
      </Button>,
      <Button disabled={selected.length === 0} key="submit" type="primary" onClick={addToPlaylist}>
        {selected.length === 0 ? 'Select Videos' : `${selected.length} videos selected`}
      </Button>,
    ];
  };

  return (
    <Modal
      title={modalTitle}
      visible={visible}
      onOk={handleOk}
      confirmLoading={confirmLoading}
      onCancel={handleCancel}
      bodyStyle={{ padding: '0px' }}
      width={1000}
      footer={modalFooter()}
    >
      <PlayistTable
        mode={mode}
        data={datastore.map((item) => { item.key = item.id; return item; })}
        setData={setDatastore}
        selected={selected}
        setSelected={setSelected}
      />
    </Modal>
  );
};

const PlayistTable = ({
  mode, data, setData, selected, setSelected,
}) => {
  switch (mode) {
    case 'EDIT':
      return <EditableTable data={data} setData={setData} />;
    default:
      return <ViewableTable data={data} selected={selected} setSelected={setSelected} />;
  }
};

const mapStateToProps = (state) => ({
  room: state.room.id,
  queue: state.room.queue,
  user: state.user.username,
});
export default connect(mapStateToProps, { createPlaylists, updatePlaylists, updateQueue })(PlaylistModel);
