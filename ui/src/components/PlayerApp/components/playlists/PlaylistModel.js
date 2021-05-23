import React, {useState, useEffect} from "react"
import { Modal, Button } from 'antd';
import {SortableTable} from "./tables/sortTable"
import { EditableTable } from "./tables/editTable";

import { Form, Input} from 'antd';

import {connect} from 'react-redux'
import {createPlaylists, updatePlaylists} from "../../../../store/playlists/playlists.actions"
import {createVideoItem, validURL} from "../../../../store/video"
import { openNotificationWithIcon } from "../../../common/notification";
const CREATE = "c"
const UPDATE = "u"

const PlaylistModel = ({visible, setVisible, data, title, room, user, createPlaylists, updatePlaylists}) => {
  const [confirmLoading, setConfirmLoading] = useState(false);
  const [sortable, setSortable] = useState(false);
  const [datastore, setDatastore] = useState([])
  const [updateType, setType] = useState(CREATE)
  const [form] = Form.useForm()


  useEffect(() => {
      if (data !== undefined){
        form.setFieldsValue({"name": data.name})
        setType(UPDATE)
        setDatastore(data.videos.map(v => {
            if(v.id === ""){
                v.id = ID()
                return v
            }
            return v 
        }))
      } else{
        form.setFieldsValue({"name": ""})
        setType(CREATE)
        setDatastore([])
      }
    }, [form, data]);
  

  const handleOk = () => {
    setConfirmLoading(true);
    savePlaylist();
    // setTimeout(() => {
    //   setVisible(false);
    //   setConfirmLoading(false);
    // }, 2000);
  };

  const handleCancel = () => {
    console.log('Clicked cancel button');
    setVisible(false);
  };
  var ID = () => {
    return '_' + Math.random().toString(36).substr(2, 9);
  };
  const addrecord = () => {
    let data = {
        id: ID(), 
        key: ID(),
        url: "",
        order: datastore.length+1,
    }
    setSortable(false)
    setDatastore(datastore => [...datastore, data]);
  };

  const savePlaylist = async () => {
      var valid = datastore.every(item => validURL(item.url));
      if (!valid){
        setConfirmLoading(false);
        openNotificationWithIcon('error', "Invalid URL")
        return;
      }
      if (form.getFieldsValue("name").name.length < 3){
        setConfirmLoading(false);
        openNotificationWithIcon('error', "Invalid: Playlist name needs to be greater then 3 characters")
        return;
      }
      
      

      let ds = await Promise.all(datastore.map(async video => { video = await createVideoItem(video.url, user); return video}))
      
      if (updateType === CREATE){
          let playlist = {
              "name": form.getFieldsValue("name").name,
              "username": "",
              "videos": ds,
              "room": room
          }
          createPlaylists(room, playlist)
      }else{
        data.name = form.getFieldsValue("name").name;
        data.videos = ds;
        updatePlaylists(room, data)
      }
    
    setVisible(false);
    setConfirmLoading(false);



  }
  console.log("playlistModel", user, room)
  return (
      <Modal
        title={title}
        visible={visible}
        onOk={handleOk}
        confirmLoading={confirmLoading}
        onCancel={handleCancel}
        bodyStyle={{padding:"0px"}}
        width={1000}
        footer={[
            <Button key="edit" onClick={()=>setSortable(!sortable)} style={{float:"left"}}>
              {sortable ? "Edit" : "Sort" }
            </Button>,
            <Button key="add" onClick={addrecord} style={{float:"left"}}>
                Add new Video
             </Button>,
            <Button key="back" onClick={handleCancel}>
              Cancel
            </Button>,
            <Button key="submit" type="primary" loading={confirmLoading} onClick={handleOk}>
              Submit
            </Button>,
        ]}
      > 
        <Form form={form}>
            <Form.Item name="name" label="Playlist Name" rules={[{ required: true }]} style={{padding: "10px 10px 0px 10px"}}>
                <Input />
            </Form.Item>
        </Form>
        {sortable ? 
        <SortableTable data={datastore.map(item => {item.key = item.id; return item})} setData={setDatastore}/> : 
        <EditableTable data={datastore.map(item => {item.key = item.id; return item})} setData={setDatastore}/> 
        }
      </Modal>
  );
};


const mapStateToProps  = (state) =>{
    return {
      room : state.room.id,
      user: state.user.username
    }
  } 
export default connect(mapStateToProps, {createPlaylists, updatePlaylists})(PlaylistModel)

