// import React, { useState,useEffect } from 'react';
// import { Table, Collapse, Button, Popconfirm, Form, Input,Space } from 'antd';
// import { listEntries } from '@/services/entry';

// const { Panel } = Collapse;

// interface DataSourceItem {
//   key: string;
//   code: string;
//   name: string;
//   amount: number;
//   weight: number;
//   note: string;
//   [key: string]: string | number; // Add an index signature
// }



// const CollapsibleTable: React.FC = () => {
//   const [expandedRowKey, setExpandedRowKey] = useState<string | undefined>(undefined);
//   const [dataSource, setDataSource] = useState<DataSourceItem[]>([]);
//   const [editingKey, setEditingKey] = useState<string | undefined>(undefined);
//   const [filters, setFilters] = useState<any>({});

//   useEffect(() => {
//     // Fetch data from your API endpoint
//     const fetchData = async () => {
//       try {
//         const response = await listEntries({
//           IsActive: true,
//           Offset: 0,
//           Limit: 9999999,
//         }); // Replace with your API endpoint
//         //setDataSource(response.data);
//         console.log(response)
//       } catch (error) {
//         console.error('Error fetching data:', error);
//       }
//     };

//     fetchData();
//   }, []); // The empty dependency array means this effect will run once after the initial render


//   const handleExpand = (key: string) => {
//     setExpandedRowKey((prevKey) => (prevKey === key ? undefined : key));
//   };

//   const handleEdit = (record: DataSourceItem) => {
//     setEditingKey(record.key);
//     setExpandedRowKey(undefined);
//   };

//   const handleDelete = (key: string) => {
//     setDataSource((prevDataSource) => prevDataSource.filter((item) => item.key !== key));
//     setEditingKey(undefined);
//     setExpandedRowKey(undefined);
//   };

//   const handleSave = (form: any, key: string) => {
//     form.validateFields((error: any, values: any) => {
//       if (error) {
//         return;
//       }
//       const newData = [...dataSource];
//       const index = newData.findIndex((item) => key === item.key);
//       if (index > -1) {
//         newData[index] = { ...newData[index], ...values };
//         setDataSource(newData);
//         setEditingKey(undefined);
//         setExpandedRowKey(undefined);
//       }
//     });
//   };

//   const columns = [
    
//     {
//       title: 'Code',
//       dataIndex: 'code',
//       key: 'code',
//     },
//     {
//       title: 'Name',
//       dataIndex: 'name',
//       key: 'name',
//     },
//     {
//       title: 'Amount',
//       dataIndex: 'amount',
//       key: 'amount',
//     },
//     {
//       title: 'Weight',
//       dataIndex: 'weight',
//       key: 'weight',
//     },
//     // {
//     //   title: 'Note',
//     //   dataIndex: 'note',
//     //   key: 'note',
//     // },
//     {
//       title: 'Actions',
//       key: 'actions',
//       render: (record:any) => (
//         <span>
//           <Button type="primary" style={{ marginRight: 8 }} onClick={() => handleEdit(record)}>
//             Edit
//           </Button>
//           <Popconfirm
//             title="Are you sure you want to delete this row?"
//             onConfirm={() => handleDelete(record.key)}
//           >
//             <Button danger>Delete</Button>
//           </Popconfirm>
//         </span>
//       ),
//     },
//   ];
//   const addRow = () => {
//     const newKey = (dataSource.length + 1).toString();
//     const newRow: DataSourceItem = {
//       key: newKey,
//       name: `New Name ${newKey}`,
//       code: `New Code ${newKey}`,
//       amount: 0,
//       weight:0.0,
//       note: 'New Note',
//     };
//     setDataSource([...dataSource, newRow]);
//     setEditingKey(newKey);
//     setExpandedRowKey(undefined);
//   };
//   const handleFilterChange = (columnKey: string, value: string) => {
//     setFilters((prevFilters:any) => ({ ...prevFilters, [columnKey]: value }));
//   };
//   const filteredDataSource = dataSource.filter((item) => {
//     return Object.keys(filters).every((key) => {
//       if (key === 'amount') {
//         return item[key] >= parseInt(filters[key], 10);
//       }
//       if (key === 'weight') {
//         return item[key] >= parseFloat(filters[key]);
//       }
//       //return item[key].toLowerCase().includes(filters[key].toLowerCase());
//       return String(item[key]).toLowerCase().includes(filters[key].toLowerCase());
//     });
//   });
//   const isEditing = (record: DataSourceItem) => record.key === editingKey;

//   const editColumns = columns.map((col) => {
//     // if (!col.editable) {
//     //   return col;
//     // }
//     const editableCol = col as {
//       title: string;
//       dataIndex: string;
//       key: string;
//       editable?: boolean; // Add editable property
//       onCell: (record: DataSourceItem) => {
//         record: DataSourceItem;
//         inputType: 'text';
//         dataIndex: string;
//         title: string;
//         editing: boolean;
//       };
//     };
  
//     return {
//       ...editableCol,
//       onCell: (record: DataSourceItem) => ({
//         record,
//         inputType: 'text',
//         dataIndex: editableCol.dataIndex,
//         title: editableCol.title,
//         editing: isEditing(record),
//       }),
//     };

//   });
//   const PanelContent: React.FC<{ record: DataSourceItem }> = ({ record }) => (
//     <div>
//       {isEditing(record) ? (
//         <Form form={form} onFinish={() => handleSave(form, record.key)}>
//           <Form.Item name="name" initialValue={record.name}>
//             <Input />
//           </Form.Item>
//           <Form.Item name="age" initialValue={record.age}>
//             <Input />
//           </Form.Item>
//           <Form.Item name="address" initialValue={record.address}>
//             <Input />
//           </Form.Item>
//           <Button type="primary" htmlType="submit">
//             Save
//           </Button>
//         </Form>
//       ) : (
//         <p>{record.name} details go here.</p>
//       )}
//     </div>
//   );

//   const [form] = Form.useForm();

//   return (
//     <div>
//       <Space style={{ marginBottom: 16 }}>
//         <Button type="primary" onClick={addRow}>
//           Add Row
//         </Button>
//         <Input
//           placeholder="Filter by Name"
//           onChange={(e) => handleFilterChange('name', e.target.value)}
//           value={filters.name || ''}
//         />
//         <Input
//           placeholder="Filter by Age"
//           onChange={(e) => handleFilterChange('age', e.target.value)}
//           value={filters.age || ''}
//         />
//       </Space>
//       <Table
//         dataSource={filteredDataSource}
//         columns={[
//           ...editColumns,
//           {
//             title: 'Details',
//             key: 'details',
//             render: (_, record) => (
//               <Collapse
//                 activeKey={expandedRowKey === record.key ? [record.key] : []}
//                 onChange={() => handleExpand(record.key)}
//               >
//                 <Panel key={record.key} header={`Details for ${record.name}`}>
//                   <PanelContent record={record} />
//                 </Panel>
//               </Collapse>
//             ),
//           },
//         ]}
//         rowKey="key"
//         pagination={{
//           pageSize: 5,
//           hideOnSinglePage: true,
//         }}
//       />
//     </div>
//   );
// };

// {/* <Tabs activeKey={activeTab} onChange={handleTabChange}>
//         <TabPane tab="Table 1" key="A">
        
//         <CollapsibleTable 
//           tabKey='A'
//         />

//         {/* </CollapsibleTable> */}
//         </TabPane>
//         <TabPane tab="Table 2" key="B">
//         <CollapsibleTable 
//             tabKey='B'
//           />
//         </TabPane>
//         </Tabs> */}

// export default CollapsibleTable;
