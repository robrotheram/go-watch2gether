import React, {
  useContext, useState, useEffect, useRef,
} from 'react';
import { Table, Input, Form } from 'antd';
import { EditableCols } from './colomns';
import { validURL } from '../../../../../store/video';

export const EditableContext = React.createContext();
const EditableRow = ({ ...props }) => {
  const [form] = Form.useForm();
  return (
    <Form form={form} component={false}>
      <EditableContext.Provider value={form}>
        <tr {...props} />
      </EditableContext.Provider>
    </Form>
  );
};
const EditableCell = ({
  title,
  editable,
  children,
  dataIndex,
  record,
  handleSave,
  ...restProps
}) => {
  const [editing, setEditing] = useState(false);
  const inputRef = useRef();

  const form = useContext(EditableContext);

  useEffect(() => {
    if (editing) {
      inputRef.current.focus();
    }
  }, [editing]);

  const toggleEdit = () => {
    setEditing(!editing);
    form.setFieldsValue({
      [dataIndex]: record[dataIndex],
    });
  };

  const save = async () => {
    try {
      const values = await form.validateFields();
      toggleEdit();
      handleSave({ ...record, ...values });
    } catch (errInfo) {
      console.log('Save failed:', errInfo);
    }
  };

  let childNode = children;

  if (editable) {
    childNode = editing ? (
      <Form.Item
        style={{
          margin: 0,
        }}
        name={dataIndex}
        rules={[
          {
            required: true,
            message: `${title} is required.`,
          },
          {
            validator: async (rule, value) => {
              if (!validURL(value)) {
                throw new Error('Invalid URL');
              }
            },
          },
        ]}
      >
        <Input ref={inputRef} onPressEnter={save} onBlur={save} />
      </Form.Item>
    ) : (
      <div
        className="editable-cell-value-wrap"
        style={{
          paddingRight: 24,
        }}
        onClick={toggleEdit}
      >
        {children}
      </div>
    );
  }

  return <td {...restProps}>{childNode}</td>;
};

export const EditableTable = ({ data, setData }) => {
  const [datastore, setDatastore] = useState(data);
  useEffect(() => { setDatastore(data); }, [data]);

  const handleDelete = (key) => {
    setData(datastore.filter((item) => item.key !== key));
  };

  const handleSave = (row) => {
    const newData = [...datastore];
    const index = newData.findIndex((item) => row.key === item.key);
    const item = newData[index];
    newData.splice(index, 1, {
      ...item,
      ...row,
    });
    setData(newData);
  };

  const columns = EditableCols(handleDelete).map((col) => {
    if (!col.editable) {
      return col;
    }
    return {
      ...col,
      onCell: (record) => ({
        record,
        editable: col.editable,
        dataIndex: col.dataIndex,
        title: col.title,
        handleSave,
      }),
    };
  });

  return (
    <Table
      style={{ height: '500px', overflowY: 'auto' }}
      components={{
        body: {
          row: EditableRow,
          cell: EditableCell,
        },
      }}
      rowKey="id"
      rowClassName={() => 'editable-row'}
      dataSource={datastore}
      columns={columns}
      pagination={false}
    />
  );
};

/*
        {/* <Button onClick={this.handleAdd} type="primary" style={{ marginBottom: 16 }}>
          Add a row
        </Button> */
