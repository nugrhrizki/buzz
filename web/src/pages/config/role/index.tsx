import { Table } from "@tanstack/solid-table";
import { RiSystemRefreshLine } from "solid-icons/ri";
import { TbLoader, TbPlus } from "solid-icons/tb";
import { Show, createEffect, createSignal } from "solid-js";

import { useRoles } from "@/services/role";

import { Role } from "@/models/role";

import { Button } from "@/components/ui/button";
import Datatable from "@/components/ui/data-table";
import { Input } from "@/components/ui/input";

import { TableDetail, columns, createRoleTable } from "./data-table";

function RolePage() {
  const role = useRoles();
  const [table, setTable] = createSignal<Table<Role>>(createRoleTable(role.data || [], columns));

  createEffect(() => {
    setTable(createRoleTable(role.data || [], columns));
  });

  return (
    <div class="space-y-4 p-8 pt-6">
      <div class="flex items-center justify-between space-y-2">
        <h2 class="text-3xl font-bold tracking-tight">
          Role
          <Show when={role.isFetching}>
            <TbLoader class="inline-flex ml-2 h-4 w-4 animate-spin" />
          </Show>
        </h2>
      </div>
      <Datatable.Root>
        <div class="flex items-end gap-x-4">
          <Input
            placeholder="Filter roles..."
            value={(table().getColumn("name")?.getFilterValue() as string) ?? ""}
            onInput={(event) => {
              table().getColumn("name")?.setFilterValue(event.target.value);
            }}
            class="max-w-sm"
          />
          <div class="ml-auto flex items-end gap-x-4">
            <Button
              variant="outline"
              size="icon"
              onClick={() => {
                role.refetch();
              }}>
              <RiSystemRefreshLine
                classList={{
                  "animate-spin": role.isFetching,
                }}
              />
            </Button>
            <Datatable.ColumnVisibility table={table()} />
            <Button>
              <TbPlus class="w-5 h-5  mr-2" />
              Tambah
            </Button>
          </div>
        </div>
        <Datatable.Table table={table()}>{(row) => <TableDetail role={row.original} />}</Datatable.Table>
        <Datatable.Pagination table={table()} />
      </Datatable.Root>
    </div>
  );
}
export default RolePage;
