import { TbPlus, TbRefresh } from "solid-icons/tb";

import { Client } from "@/models/client";

import { Button } from "@/components/ui/button";
import Datatable from "@/components/ui/data-table";
import { Input } from "@/components/ui/input";

import { TableDetail, columns, createClientTable } from "./data-table";

const data: Client[] = [
  {
    id: 1,
    name: "Admin",
    token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
    connected: 1,
    jid: "628234422522",
  },
  {
    id: 2,
    name: "User",
    token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
    connected: 1,
    jid: "628234422522",
  },
];

function ClientPage() {
  const table = createClientTable(data, columns);

  return (
    <div class="space-y-4 p-8 pt-6">
      <div class="flex items-center justify-between space-y-2">
        <h2 class="text-3xl font-bold tracking-tight">Client</h2>
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
            <Button variant="outline" size="icon">
              <TbRefresh />
            </Button>
            <Datatable.ColumnVisibility table={table} />
            <Button>
              <TbPlus class="w-5 h-5  mr-2" />
              Tambah
            </Button>
          </div>
        </div>
        <Datatable.Table table={table}>{(row) => <TableDetail client={row.original} />}</Datatable.Table>
        <Datatable.Pagination table={table} />
      </Datatable.Root>
    </div>
  );
}
export default ClientPage;
