import React, {useState, useEffect} from "react"
import { Modal, Button } from 'antd';
import {SortableTable} from "./tables/sortTable"
import { EditableTable } from "./tables/editTable";

import { Form, Input} from 'antd';
import { FormInstance } from 'antd/lib/form';

import {connect} from 'react-redux'
import {createPlaylists, updatePlaylists} from "../../../../store/playlists/playlists.actions"

const CREATE = "c"
const UPDATE = "u"
const initdata = [];

const PlaylistModel = ({visible, setVisible, data, title, room, createPlaylists, updatePlaylists}) => {
  const [confirmLoading, setConfirmLoading] = useState(false);
  const [sortable, setSortable] = useState(false);
  const [modalText, setModalText] = useState('Content of the modal');
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
    }, [data]);
  const showModal = () => {
    setVisible(true);
  };

  const handleOk = () => {
    setModalText('The modal will be closed after two seconds');
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
        key: ID(),
        url: "",
        order: datastore.length+1,
    }
    setSortable(false)
    setDatastore(datastore => [...datastore, data]);
  };

  const savePlaylist = () => {
      if (updateType === CREATE){
          let playlist = {
              "name": form.getFieldsValue("name").name,
              "username": "",
              "videos": datastore,
              "room": room
          }
          createPlaylists(room, playlist)
      }else{
        data.name = form.getFieldsValue("name").name;
        data.videos = datastore;
        updatePlaylists(room, data)
      }
    
    setVisible(false);
    setConfirmLoading(false);



  }

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
        <SortableTable data={datastore} setData={setDatastore}/> : 
        <EditableTable  data={datastore} setData={setDatastore}/> 
        }
      </Modal>
  );
};


const mapStateToProps  = (state) =>{
    return {
      room : state.room.id,
      playlists : state.playlist.playlists,
      loading: state.playlist.loading
    }
  } 
export default connect(mapStateToProps, {createPlaylists, updatePlaylists})(PlaylistModel)

