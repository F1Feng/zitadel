import { SelectionModel } from '@angular/cdk/collections';
import { Component, EventEmitter, Input, OnInit, Output, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatTableDataSource } from '@angular/material/table';
import { ActivatedRoute, Router } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { take } from 'rxjs/operators';
import { enterAnimations } from 'src/app/animations';
import { ActionKeysType } from 'src/app/modules/action-keys/action-keys.component';
import { PageEvent, PaginatorComponent } from 'src/app/modules/paginator/paginator.component';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import { Timestamp } from 'src/app/proto/generated/google/protobuf/timestamp_pb';
import { SearchQuery, Type, TypeQuery, User, UserState } from 'src/app/proto/generated/zitadel/user_pb';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

enum UserListSearchKey {
  FIRST_NAME,
  LAST_NAME,
  DISPLAY_NAME,
  USER_NAME,
  EMAIL,
}

@Component({
  selector: 'cnsl-user-table',
  templateUrl: './user-table.component.html',
  styleUrls: ['./user-table.component.scss'],
  animations: [enterAnimations],
})
export class UserTableComponent implements OnInit {
  public userSearchKey: UserListSearchKey | undefined = undefined;
  public Type: any = Type;
  @Input() public type: Type = Type.TYPE_HUMAN;
  @Input() refreshOnPreviousRoutes: string[] = [];
  @Input() disabled: boolean = false;
  @ViewChild(PaginatorComponent) public paginator!: PaginatorComponent;
  @ViewChild('input') public filter!: Input;

  public viewTimestamp!: Timestamp.AsObject;
  public totalResult: number = 0;
  public dataSource: MatTableDataSource<User.AsObject> = new MatTableDataSource<User.AsObject>();
  public selection: SelectionModel<User.AsObject> = new SelectionModel<User.AsObject>(true, []);
  private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public loading$: Observable<boolean> = this.loadingSubject.asObservable();
  @Input() public displayedColumnsHuman: string[] = ['select', 'displayName', 'username', 'email', 'state', 'actions'];
  @Input() public displayedColumnsMachine: string[] = ['select', 'displayName', 'username', 'state', 'actions'];

  @Output() public changedSelection: EventEmitter<Array<User.AsObject>> = new EventEmitter();

  public UserState: any = UserState;
  public UserListSearchKey: any = UserListSearchKey;

  public ActionKeysType: any = ActionKeysType;
  public filterOpen: boolean = false;

  constructor(
    private router: Router,
    public translate: TranslateService,
    private authService: GrpcAuthService,
    private userService: ManagementService,
    private toast: ToastService,
    private dialog: MatDialog,
    private route: ActivatedRoute,
  ) {
    this.selection.changed.subscribe(() => {
      this.changedSelection.emit(this.selection.selected);
    });
  }

  ngOnInit(): void {
    this.route.queryParams.pipe(take(1)).subscribe((params) => {
      this.getData(10, 0, this.type);
      if (params.deferredReload) {
        setTimeout(() => {
          this.getData(this.paginator.pageSize, this.paginator.pageIndex * this.paginator.pageSize, this.type);
        }, 2000);
      }
    });
  }

  public setType(type: Type): void {
    this.type = type;
    this.router.navigate([], {
      relativeTo: this.route,
      queryParams: {
        type: type === Type.TYPE_HUMAN ? 'human' : type === Type.TYPE_MACHINE ? 'machine' : 'human',
      },
      queryParamsHandling: 'merge',
      skipLocationChange: false,
    });
    this.getData(this.paginator.pageSize, this.paginator.pageIndex * this.paginator.pageSize, this.type);
  }

  public isAllSelected(): boolean {
    const numSelected = this.selection.selected.length;
    const numRows = this.dataSource.data.length;
    return numSelected === numRows;
  }

  public masterToggle(): void {
    this.isAllSelected() ? this.selection.clear() : this.dataSource.data.forEach((row) => this.selection.select(row));
  }

  public changePage(event: PageEvent): void {
    this.selection.clear();
    this.getData(event.pageSize, event.pageIndex * event.pageSize, this.type);
  }

  public deactivateSelectedUsers(): void {
    Promise.all(
      this.selection.selected
        .filter((u) => u.state === UserState.USER_STATE_ACTIVE)
        .map((value) => {
          return this.userService.deactivateUser(value.id);
        }),
    )
      .then(() => {
        this.toast.showInfo('USER.TOAST.SELECTEDDEACTIVATED', true);
        this.selection.clear();
        setTimeout(() => {
          this.refreshPage();
        }, 1000);
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public reactivateSelectedUsers(): void {
    Promise.all(
      this.selection.selected
        .filter((u) => u.state === UserState.USER_STATE_INACTIVE)
        .map((value) => {
          return this.userService.reactivateUser(value.id);
        }),
    )
      .then(() => {
        this.toast.showInfo('USER.TOAST.SELECTEDREACTIVATED', true);
        this.selection.clear();
        setTimeout(() => {
          this.refreshPage();
        }, 1000);
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public gotoRouterLink(rL: any): void {
    this.router.navigate(rL);
  }

  private async getData(limit: number, offset: number, type: Type, searchQueries?: SearchQuery[]): Promise<void> {
    this.loadingSubject.next(true);

    let queryT = new SearchQuery();
    const typeQuery = new TypeQuery();
    typeQuery.setType(type);
    queryT.setTypeQuery(typeQuery);

    this.userService
      .listUsers(limit, offset, searchQueries?.length ? [queryT, ...searchQueries] : [queryT])
      .then((resp) => {
        if (resp.details?.totalResult) {
          this.totalResult = resp.details?.totalResult;
        } else {
          this.totalResult = 0;
        }
        if (resp.details?.viewTimestamp) {
          this.viewTimestamp = resp.details?.viewTimestamp;
        }
        this.dataSource.data = resp.resultList;
        this.loadingSubject.next(false);
      })
      .catch((error) => {
        this.toast.showError(error);
        this.loadingSubject.next(false);
      });
  }

  public refreshPage(): void {
    this.getData(this.paginator.pageSize, this.paginator.pageIndex * this.paginator.pageSize, this.type);
  }

  public applySearchQuery(searchQueries: SearchQuery[]): void {
    this.selection.clear();
    this.getData(this.paginator.pageSize, this.paginator.pageIndex * this.paginator.pageSize, this.type, searchQueries);
  }

  public setFilter(key: UserListSearchKey): void {
    setTimeout(() => {
      if (this.filter) {
        (this.filter as any).nativeElement.focus();
      }
    }, 100);

    if (this.userSearchKey !== key) {
      this.userSearchKey = key;
    } else {
      this.userSearchKey = undefined;
      this.refreshPage();
    }
  }

  public deleteUser(user: User.AsObject): void {
    const authUserData = {
      confirmKey: 'ACTIONS.DELETE',
      cancelKey: 'ACTIONS.CANCEL',
      titleKey: 'USER.DIALOG.DELETE_SELF_TITLE',
      descriptionKey: 'USER.DIALOG.DELETE_SELF_DESCRIPTION',
      confirmationKey: 'USER.DIALOG.TYPEUSERNAME',
      confirmation: user.preferredLoginName,
    };

    const mgmtUserData = {
      confirmKey: 'ACTIONS.DELETE',
      cancelKey: 'ACTIONS.CANCEL',
      titleKey: 'USER.DIALOG.DELETE_TITLE',
      descriptionKey: 'USER.DIALOG.DELETE_DESCRIPTION',
      confirmationKey: 'USER.DIALOG.TYPEUSERNAME',
      confirmation: user.preferredLoginName,
    };

    if (user && user.id) {
      const authUser = this.authService.userSubject.getValue();
      const isMe = authUser?.id === user.id;

      let dialogRef;

      if (isMe) {
        dialogRef = this.dialog.open(WarnDialogComponent, {
          data: authUserData,
          width: '400px',
        });
      } else {
        dialogRef = this.dialog.open(WarnDialogComponent, {
          data: mgmtUserData,
          width: '400px',
        });
      }

      dialogRef.afterClosed().subscribe((resp) => {
        if (resp) {
          this.userService
            .removeUser(user.id)
            .then(() => {
              setTimeout(() => {
                this.refreshPage();
              }, 1000);
              this.toast.showInfo('USER.TOAST.DELETED', true);
            })
            .catch((error) => {
              this.toast.showError(error);
            });
        }
      });
    }
  }
}
