import React, { useContext, useState, useEffect, useRef } from 'react';
import { Table, Input, Form } from 'antd';
import {cols} from "./colomns"

export const ViewableTable = ({data, selected, setSelected, }) => {
  const [datastore, setDatastore] = useState(data)
  useEffect(() => {setDatastore(data)}, [data]);


  const selectRow = (record) => {
    if (selected.indexOf(record.key) >= 0) {
      selected.splice(selected.indexOf(record.key), 1);
    } else {
      selected.push(record.key);
    }
    setSelected(selected)
  }

  const onSelectedRowKeysChange = (selected) => {
    setSelected(selected)
  }

  const rowSelection = {
    selected,
    columnWidth: 80,
    onChange: onSelectedRowKeysChange,
  };
    return (
      <Table
        rowSelection={rowSelection}
        style={{height:"500px", "overflowY":"auto"}}
        columns={cols}
        dataSource={datastore}
        pagination={false}
        onRow={(record) => ({
          onClick: () => {
            selectRow(record);
          },
        })}
      />
    );
}

/*
        {/* <Button onClick={this.handleAdd} type="primary" style={{ marginBottom: 16 }}>
          Add a row
        </Button> */
        