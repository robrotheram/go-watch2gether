import React, {useState}from 'react'
import {connect} from 'react-redux'
import { Drawer, Space, Button, Col, Row, Input, Select, DatePicker, List } from 'antd';
import { PlusOutlined,PlaySquareOutlined, DeleteOutlined, EditOutlined, SelectOutlined } from '@ant-design/icons';
import UserList from '../UserList';
import Share from '../ShareModal'
import Settings from '../SettingsModal'
import {
  SyncOutlined
} from '@ant-design/icons';

import { Popconfirm } from 'antd';

import {getPlaylists, deletePlaylists, loadPlaylists} from "../../../../store/playlists/playlists.actions"
import {VideoItem} from "../VideoQueue/VideoItem"

import { PlaylistItem } from './PlaylistItem';
import PlaylistModel from './PlaylistModel';

const { Option } = Select;

const PlaylistDrawer = (props) => {

  const [visible, setVisible] = useState(false);

  const [modalVisible, setModalVisible] = useState(false);
  const [modelTitle, setModelTitle] = useState("Create new Playlist");
  const [datastore, setDatastore] = useState()
  
  const showDrawer = () => { 
    console.log(props.room)
    props.getPlaylists(props.room)
    setVisible(true); 
  };
  const onClose = () => { setVisible(false); };

  const showModel = (playlist) => {
    console.log("SHOWMODLE", playlist)
    if (playlist === undefined) {
      setModelTitle("Create new Playlist")
      setDatastore()
    }else{
      setModelTitle("Edit playlist: "+playlist.name)
      setDatastore(playlist)
    }
    setModalVisible(true)
  }

  const deletePlaylist = (playlist) => {
    props.deletePlaylists(playlist.room, playlist.id)
  }

  const loadPlaylist = (playlist) => {
    props.loadPlaylists(playlist.room, playlist.id)
  }


    return (
      <>
        <Button type="primary" onClick={showDrawer} style={{"height":"33px", "margin": "0px 10px"}} >
          <PlaySquareOutlined />Playlists
        </Button>
        <Drawer
          title={
          <Space>
            <Button type="primary" onClick={()=>showModel()}><PlusOutlined /></Button>
            Playlists
          </Space>
          }
          width={470}
          onClose={onClose}
          visible={visible}
          maskClosable={true}
          placement="left"
          bodyStyle={{ padding: 0 }}
        >
        <div style={{height:"100%", width:"100%", overflowX:"hidden"}}>
        <List
                size="small"
                itemLayout="horizontal"
                dataSource={props.playlists} //deletePlaylist
                renderItem={item => (
                  <PlaylistItem key={item.id} video={item.videos[0]} playlist={item} loading={item.loading}>
                      <Space>
                                <Popconfirm style={{width:"90px"}} title="Sure to delete?" onConfirm={() => deletePlaylist(item)}>
                                <Button style={{width:"90px"}} icon={<DeleteOutlined/>}>Delete</Button> 
                                </Popconfirm>
                                <Button style={{width:"90px"}} icon={<EditOutlined/>} onClick={()=>showModel(item)} >Edit</Button>
                                <Popconfirm style={{width:"90px"}} title="Sure to Load Playlist?" onConfirm={() => loadPlaylist(item)}>
                                <Button style={{width:"90px"}} icon={<SelectOutlined/>}>Load</Button> 
                                </Popconfirm>
                                                     
                      </Space>
                  </PlaylistItem>
                )}
            />
          </div>
          <PlaylistModel visible={modalVisible} setVisible={setModalVisible} title={modelTitle} data={datastore}/>
        </Drawer>
      </>
    );
}
const mapStateToProps  = (state) =>{
  return {
    room : state.room.id,
    playlists : state.playlist.playlists,
    loading: state.playlist.loading
  }
} 
export default connect(mapStateToProps, {getPlaylists, deletePlaylists, loadPlaylists})(PlaylistDrawer)
