import { TbPlus } from "solid-icons/tb";

import { Role } from "@/models/role";

import { Button } from "@/components/ui/button";
import Datatable from "@/components/ui/data-table";
import { Input } from "@/components/ui/input";

import { TableDetail, columns, createRoleTable } from "./data-table";

const data: Role[] = [
  {
    id: 1,
    name: "Admin",
    actions: "{}",
    created_at: new Date("2021-08-31T08:00:00.000Z"),
    updated_at: new Date("2021-08-31T08:00:00.000Z"),
    deleted_at: null,
  },
  {
    id: 2,
    name: "User",
    actions: "{}",
    created_at: new Date("2021-08-31T08:00:00.000Z"),
    updated_at: new Date("2021-08-31T08:00:00.000Z"),
    deleted_at: null,
  },
  {
    id: 3,
    name: "Guest",
    actions: "{}",
    created_at: new Date("2021-08-31T08:00:00.000Z"),
    updated_at: new Date("2021-08-31T08:00:00.000Z"),
    deleted_at: null,
  },
];

function RolePage() {
  const table = createRoleTable(data, columns);

  return (
    <div class="space-y-4 p-8 pt-6">
      <div class="flex items-center justify-between space-y-2">
        <h2 class="text-3xl font-bold tracking-tight">Role</h2>
      </div>
      <Datatable.Root>
        <div class="flex items-end gap-x-4">
          <Input
            placeholder="Filter roles..."
            value={(table.getColumn("name")?.getFilterValue() as string) ?? ""}
            onInput={(event) => {
              table.getColumn("name")?.setFilterValue(event.target.value);
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
        <Datatable.Table table={table}>{(row) => <TableDetail role={row.original} />}</Datatable.Table>
        <Datatable.Pagination table={table} />
      </Datatable.Root>
    </div>
  );
}
export default RolePage;
