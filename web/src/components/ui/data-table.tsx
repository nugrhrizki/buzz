import { As } from "@kobalte/core";
import { Column, ColumnDef, Row, Table as TanstackTable, flexRender } from "@tanstack/solid-table";
import {
  TbAdjustmentsHorizontal,
  TbArrowDown,
  TbArrowUp,
  TbArrowsUpDown,
  TbChevronDown,
  TbChevronLeft,
  TbChevronRight,
  TbChevronsLeft,
  TbChevronsRight,
  TbEyeOff,
} from "solid-icons/tb";
import { For, JSX, Match, ParentProps, Show, Switch } from "solid-js";

import { cn } from "@/lib/utils";

import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";

function Root(props: ParentProps) {
  return props.children;
}

interface DataTableProps<T> {
  table: TanstackTable<T>;
  children?: (row: Row<T>) => JSX.Element;
}

function DataTable<T>(props: DataTableProps<T>) {
  function columnsLength() {
    let length = 0;
    for (const header of props.table.getHeaderGroups()) {
      if (header.headers.length > length) {
        length = header.headers.length;
      }
    }

    return length;
  }

  return (
    <div class="w-full">
      <div class="rounded-md border">
        <Table>
          <TableHeader>
            <For each={props.table.getHeaderGroups()}>
              {(headerGroup) => (
                <TableRow>
                  <For each={headerGroup.headers}>
                    {(header) => (
                      <TableHead>
                        <Show when={!header.isPlaceholder}>
                          {flexRender(header.column.columnDef.header, header.getContext())}
                        </Show>
                      </TableHead>
                    )}
                  </For>
                </TableRow>
              )}
            </For>
          </TableHeader>
          <TableBody>
            <For
              each={props.table.getRowModel().rows}
              fallback={
                <TableRow>
                  <TableCell colSpan={columnsLength()} class="h-24 text-center">
                    No results.
                  </TableCell>
                </TableRow>
              }>
              {(row) => (
                <>
                  <TableRow data-state={row.getIsSelected() && "selected"}>
                    <For each={row.getVisibleCells()}>
                      {(cell) => {
                        return (
                          <TableCell
                            class={
                              (cell.column.columnDef.id === "expander" && row.getCanExpand()) ||
                              (cell.column.columnDef.id === "select" && row.getCanSelect())
                                ? "w-[1%] whitespace-nowrap"
                                : ""
                            }>
                            {flexRender(cell.column.columnDef.cell, cell.getContext())}
                          </TableCell>
                        );
                      }}
                    </For>
                  </TableRow>
                  <Show when={row.getIsExpanded() && props.children}>
                    <TableRow>
                      <TableCell colSpan={row.getVisibleCells().length} class="p-0">
                        {props.children!(row)}
                      </TableCell>
                    </TableRow>
                  </Show>
                </>
              )}
            </For>
          </TableBody>
        </Table>
      </div>
    </div>
  );
}

interface DataTableViewOptionsProps<TData> {
  table: TanstackTable<TData>;
}

function DataTableViewOptions<TData>(props: DataTableViewOptionsProps<TData>) {
  function columns() {
    return props.table
      .getAllColumns()
      .filter((column) => typeof column.accessorFn !== "undefined" && column.getCanHide());
  }

  return (
    <Show when={columns().length > 0}>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <As component={Button} variant="outline">
            <TbAdjustmentsHorizontal class="mr-2 h-4 w-4" />
            Columns
          </As>
        </DropdownMenuTrigger>
        <DropdownMenuContent class="w-[150px]">
          <For each={columns()}>
            {(column) => (
              <DropdownMenuCheckboxItem
                class="capitalize"
                checked={column.getIsVisible()}
                onChange={(value) => column.toggleVisibility(!!value)}>
                {column.id}
              </DropdownMenuCheckboxItem>
            )}
          </For>
        </DropdownMenuContent>
      </DropdownMenu>
    </Show>
  );
}

function DatatableExpandableRow<T>(): ColumnDef<T> {
  return {
    id: "expander",
    header: () => null,
    cell: (cell) => (
      <Show when={cell.row.getCanExpand()}>
        <Button
          variant="ghost"
          size="sm"
          class="text-muted-foreground cursor-pointer"
          onClick={cell.row.getToggleExpandedHandler()}>
          <Show when={cell.row.getIsExpanded()} fallback={<TbChevronRight class="w-4 h-4" />}>
            <TbChevronDown class="w-4 h-4" />
          </Show>
        </Button>
      </Show>
    ),
    enableHiding: false,
  };
}

function DatatableSelectableRow<T>(): ColumnDef<T> {
  return {
    id: "select",
    header: (header) => (
      <Checkbox
        checked={header.table.getIsAllPageRowsSelected()}
        onChange={(value) => header.table.toggleAllPageRowsSelected(!!value)}
        aria-label="Select all"
      />
    ),
    cell: (cell) => (
      <Checkbox
        checked={cell.row.getIsSelected()}
        onChange={(value) => cell.row.toggleSelected(!!value)}
        aria-label="Select row"
      />
    ),
    enableSorting: false,
    enableHiding: false,
  };
}

interface DataTableColumnHeaderProps<TData, TValue> extends JSX.HTMLAttributes<HTMLDivElement> {
  column: Column<TData, TValue>;
  title: string;
}

function DataTableColumnHeader<TData, TValue>(props: DataTableColumnHeaderProps<TData, TValue>) {
  if (!props.column.getCanSort() && !props.column.getCanHide()) {
    return <div class={cn(props.class)}>{props.title}</div>;
  }

  return (
    <div class={cn("flex items-center space-x-2", props.class)}>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <As component={Button} variant="ghost" size="sm" class="-ml-3 h-8 data-[state=open]:bg-accent">
            <span>{props.title}</span>
            <Show when={props.column.getCanSort()}>
              <Switch fallback={<TbArrowsUpDown class="ml-2 h-4 w-4" />}>
                <Match when={props.column.getIsSorted() === "desc"}>
                  <TbArrowDown class="ml-2 h-4 w-4" />
                </Match>
                <Match when={props.column.getIsSorted() === "asc"}>
                  <TbArrowUp class="ml-2 h-4 w-4" />
                </Match>
              </Switch>
            </Show>
          </As>
        </DropdownMenuTrigger>
        <DropdownMenuContent>
          <Show when={props.column.getCanSort()}>
            <DropdownMenuItem onClick={() => props.column.toggleSorting(false)}>
              <TbArrowUp class="mr-2 h-3.5 w-3.5 text-muted-foreground/70" />
              Asc
            </DropdownMenuItem>
            <DropdownMenuItem onClick={() => props.column.toggleSorting(true)}>
              <TbArrowDown class="mr-2 h-3.5 w-3.5 text-muted-foreground/70" />
              Desc
            </DropdownMenuItem>
          </Show>
          <Show when={props.column.getCanHide()}>
            <DropdownMenuSeparator />
            <DropdownMenuItem onClick={() => props.column.toggleVisibility(false)}>
              <TbEyeOff class="mr-2 h-3.5 w-3.5 text-muted-foreground/70" />
              Hide
            </DropdownMenuItem>
          </Show>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  );
}

interface DataTablePaginationProps<TData> {
  table: TanstackTable<TData>;
}

function DataTablePagination<TData>(props: DataTablePaginationProps<TData>) {
  return (
    <div class="flex items-center justify-between px-2">
      <div class="flex-1 text-sm text-muted-foreground">
        {props.table.getFilteredSelectedRowModel().rows.length} of {props.table.getFilteredRowModel().rows.length}{" "}
        row(s) selected.
      </div>
      <div class="flex items-center space-x-6 lg:space-x-8">
        <div class="flex items-center space-x-2">
          <p class="text-sm font-medium">Rows per page</p>
          <Select
            value={`${props.table.getState().pagination.pageSize}`}
            onChange={(value) => {
              props.table.setPageSize(Number(value));
            }}
            options={[10, 20, 30, 40, 50]}
            placeholder={props.table.getState().pagination.pageSize}
            itemComponent={(props) => <SelectItem item={props.item}>{props.item.rawValue}</SelectItem>}>
            <SelectTrigger class="h-8 w-[70px]">
              <SelectValue<string>>{(state) => state.selectedOption()}</SelectValue>
            </SelectTrigger>
            <SelectContent />
          </Select>
        </div>
        <div class="flex w-[100px] items-center justify-center text-sm font-medium">
          Page {props.table.getState().pagination.pageIndex + 1} of {props.table.getPageCount()}
        </div>
        <div class="flex items-center space-x-2">
          <Button
            variant="outline"
            class="hidden h-8 w-8 p-0 lg:flex"
            onClick={() => props.table.setPageIndex(0)}
            disabled={!props.table.getCanPreviousPage()}>
            <span class="sr-only">Go to first page</span>
            <TbChevronsLeft class="h-4 w-4" />
          </Button>
          <Button
            variant="outline"
            class="h-8 w-8 p-0"
            onClick={() => props.table.previousPage()}
            disabled={!props.table.getCanPreviousPage()}>
            <span class="sr-only">Go to previous page</span>
            <TbChevronLeft class="h-4 w-4" />
          </Button>
          <Button
            variant="outline"
            class="h-8 w-8 p-0"
            onClick={() => props.table.nextPage()}
            disabled={!props.table.getCanNextPage()}>
            <span class="sr-only">Go to next page</span>
            <TbChevronRight class="h-4 w-4" />
          </Button>
          <Button
            variant="outline"
            class="hidden h-8 w-8 p-0 lg:flex"
            onClick={() => props.table.setPageIndex(props.table.getPageCount() - 1)}
            disabled={!props.table.getCanNextPage()}>
            <span class="sr-only">Go to last page</span>
            <TbChevronsRight class="h-4 w-4" />
          </Button>
        </div>
      </div>
    </div>
  );
}

export default {
  Root,
  Table: DataTable,
  Pagination: DataTablePagination,
  ColumnHeader: DataTableColumnHeader,
  ColumnVisibility: DataTableViewOptions,
  RowExpand: DatatableExpandableRow,
  RowSelect: DatatableSelectableRow,
};
