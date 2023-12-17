import { TbPlus } from "solid-icons/tb";

import { User } from "@/models/user";

import { Button } from "@/components/ui/button";
import Datatable from "@/components/ui/data-table";
import { Input } from "@/components/ui/input";

import { TableDetail, columns, createUserTable } from "./data-table";

const data: User[] = [
  {
    id: 1,
    name: "Admin",
    username: "admin",
    email: "",
    whatsapp: "",
    confirmed: true,
    password: "",
    role_id: 1,
    created_at: new Date("2021-08-31T08:00:00.000Z"),
    updated_at: new Date("2021-08-31T08:00:00.000Z"),
    deleted_at: null,
  },
  {
    id: 2,
    name: "User",
    username: "user",
    email: "",
    whatsapp: "",
    confirmed: true,
    password: "",
    role_id: 2,
    created_at: new Date("2021-08-31T08:00:00.000Z"),
    updated_at: new Date("2021-08-31T08:00:00.000Z"),
    deleted_at: null,
  },
  {
    id: 3,
    name: "Guest",
    username: "guest",
    email: "",
    whatsapp: "",
    confirmed: true,
    password: "",
    role_id: 3,
    created_at: new Date("2021-08-31T08:00:00.000Z"),
    updated_at: new Date("2021-08-31T08:00:00.000Z"),
    deleted_at: null,
  },
  {
    id: 4,
    name: "Guest",
    username: "guest",
    email: "",
    whatsapp: "",
    confirmed: true,
    password: "",
    role_id: 3,
    created_at: new Date("2021-08-31T08:00:00.000Z"),
    updated_at: new Date("2021-08-31T08:00:00.000Z"),
    deleted_at: null,
  },
  {
    id: 5,
    name: "Guest",
    username: "guest",
    email: "",
    whatsapp: "",
    confirmed: true,
    password: "",
    role_id: 3,
    created_at: new Date("2021-08-31T08:00:00.000Z"),
    updated_at: new Date("2021-08-31T08:00:00.000Z"),
    deleted_at: null,
  },
  {
    id: 6,
    name: "Guest",
    username: "guest",
    email: "",
    whatsapp: "",
    confirmed: true,
    password: "",
    role_id: 3,
    created_at: new Date("2021-08-31T08:00:00.000Z"),
    updated_at: new Date("2021-08-31T08:00:00.000Z"),
    deleted_at: null,
  },
  {
    id: 7,
    name: "Guest",
    username: "guest",
    email: "",
    whatsapp: "",
    confirmed: true,
    password: "",
    role_id: 3,
    created_at: new Date("2021-08-31T08:00:00.000Z"),
    updated_at: new Date("2021-08-31T08:00:00.000Z"),
    deleted_at: null,
  },
  {
    id: 8,
    name: "Guest",
    username: "guest",
    email: "",
    whatsapp: "",
    confirmed: true,
    password: "",
    role_id: 3,
    created_at: new Date("2021-08-31T08:00:00.000Z"),
    updated_at: new Date("2021-08-31T08:00:00.000Z"),
    deleted_at: null,
  },
];

function UserPage() {
  const table = createUserTable(data, columns);

  return (
    <div class="space-y-4 p-8 pt-6">
      <div class="flex items-center justify-between space-y-2">
        <h2 class="text-3xl font-bold tracking-tight">User</h2>
      </div>
      <Datatable.Root>
        <div class="flex items-end gap-x-4">
          <Input
            placeholder="Filter username..."
            value={(table.getColumn("username")?.getFilterValue() as string) ?? ""}
            onInput={(event) => {
              table.getColumn("username")?.setFilterValue(event.target.value);
            }}
            class="max-w-sm"
          />
          <div class="ml-auto flex items-end gap-x-4">
            <Datatable.ColumnVisibility table={table} />
            <Button>
              <TbPlus class="w-5 h-5  mr-2" />
              Tambah
            </Button>
          </div>
        </div>
        <Datatable.Table table={table}>{(row) => <TableDetail user={row.original} />}</Datatable.Table>
        <Datatable.Pagination table={table} />
      </Datatable.Root>
    </div>
  );
}
export default UserPage;
