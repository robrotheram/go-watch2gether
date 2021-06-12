import React, { useState } from 'react';
import { connect } from 'react-redux';
import {
  Drawer, Space, Button, List, Popconfirm,
} from 'antd';
import {
  PlusOutlined, PlaySquareOutlined, DeleteOutlined
} from '@ant-design/icons';

import { getPlaylists, deletePlaylists, loadPlaylists } from '../../../../store/playlists/playlists.actions';

import { PlaylistItem } from './PlaylistItem';
import PlaylistModel from './PlaylistModel';
import useDeviceDetect from '../../../common/useDeviceDetect';

const PlaylistDrawer = (props) => {
  const isMoble = useDeviceDetect();
  const drawerWidth = isMoble ? '100%' : '470px';

  const [visible, setVisible] = useState(false);

  const [modalVisible, setModalVisible] = useState(false);
  const [modelTitle, setModelTitle] = useState('Create new Playlist');
  const [datastore, setDatastore] = useState();

  const showDrawer = () => {
    console.log(props.room);
    props.getPlaylists(props.room);
    setVisible(true);
  };
  const onClose = () => { setVisible(false); };

  const showModel = (playlist) => {
    console.log('SHOWMODLE', playlist);
    if (playlist === undefined) {
      setModelTitle('Create new Playlist');
      setDatastore();
    } else {
      setModelTitle(playlist.name);
      setDatastore(playlist);
    }
    setModalVisible(true);
  };

  const deletePlaylist = (playlist) => {
    props.deletePlaylists(playlist.room, playlist.id);
  };

  return (
    <>
      <Button type="primary" onClick={showDrawer} style={{ height: '33px', margin: '0px 10px' }}>
        <PlaySquareOutlined />
        Playlists
      </Button>
      <Drawer
        title={(
          <Space>
            <Button type="primary" onClick={() => showModel()}><PlusOutlined /></Button>
            Playlists
          </Space>
        )}
        width={drawerWidth}
        onClose={onClose}
        visible={visible}
        maskClosable
        placement="left"
        bodyStyle={{ padding: 0 }}
      >
        <div style={{ height: '100%', width: '100%', overflowX: 'hidden' }}>
          <List
            size="small"
            itemLayout="horizontal"
            dataSource={props.playlists} // deletePlaylist
            renderItem={(item) => (
              <PlaylistItem key={item.id} video={item.videos[0]} playlist={item} loading={item.loading} click={() => showModel(item)}>
                <Space>
                  <Popconfirm
                    onClick={(e) => { e.stopPropagation(); }}
                    style={{ width: '90px' }}
                    title="Sure to delete?"
                    onConfirm={(e) => { e.stopPropagation(); deletePlaylist(item); }}
                    onCancel={(e) => { e.stopPropagation(); }}
                  >
                    <Button icon={<DeleteOutlined />} />
                  </Popconfirm>
                  {/* <Button style={{width:"90px"}} icon={<EditOutlined/>} onClick={()=>showModel(item)} >Edit</Button> */}

                </Space>
              </PlaylistItem>
            )}
          />
        </div>
        <PlaylistModel visible={modalVisible} setVisible={setModalVisible} title={modelTitle} data={datastore} />
      </Drawer>
    </>
  );
};
const mapStateToProps = (state) => ({
  room: state.room.id,
  playlists: state.playlist.playlists,
  loading: state.playlist.loading,
});
export default connect(mapStateToProps, { getPlaylists, deletePlaylists, loadPlaylists })(PlaylistDrawer);
